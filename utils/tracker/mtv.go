package tracker

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/imroc/req"
	"github.com/jpillora/backoff"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gitlab.com/cloudb0x/trackarr/config"
	"gitlab.com/cloudb0x/trackarr/utils/maps"
	"gitlab.com/cloudb0x/trackarr/utils/web"
	"go.uber.org/ratelimit"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

/* Const */
const (
	mtvLoginUrl         = "https://www.morethan.tv/login.php"
	mtvTorrentUrl       = "https://www.morethan.tv/ajax.php"
	mtvTimeout          = 30
	mtvApiRateLimit     = 1
	mtvMaxLoginAttempts = 5
)

/* Var */
var (
	// errors
	errMtvMaxLoginAttempts = errors.New("max login attempts reached")
)

/* Struct */
type Mtv struct {
	log           *logrus.Entry
	tracker       *config.TrackerInstance
	loginBody     req.Param
	rl            *ratelimit.Limiter
	mtx           sync.Mutex
	cookieExpiry  *time.Time
	loginAttempts int
}

/* Private */

func newMtv(tracker *config.TrackerInstance) (Interface, error) {
	log := log.WithField("api", tracker.Name)

	// validate required tracker settings available
	siteUser, err := maps.GetStringMapValue(tracker.Config.Settings, "site_user", false)
	if err != nil {
		return nil, errors.WithMessage(err, "site_user setting missing")
	}

	sitePass, err := maps.GetStringMapValue(tracker.Config.Settings, "site_pass", false)
	if err != nil {
		return nil, errors.WithMessage(err, "site_pass setting missing")
	}

	// return api instance
	return &Mtv{
		log:     log,
		tracker: tracker,
		loginBody: req.Param{
			"username":   siteUser,
			"password":   sitePass,
			"login":      "Log in",
			"keeplogged": 1,
		},
		rl:            web.GetRateLimiter(tracker.Name, mtvApiRateLimit),
		loginAttempts: 0,
	}, nil
}

func (t *Mtv) parseTorrentId(torrentUrl string) (string, error) {
	// parse url
	u, err := url.Parse(torrentUrl)
	if err != nil {
		return "", errors.Wrap(err, "failed parsing torrent url")
	}

	// get id
	if id := u.Query().Get("id"); id != "" {
		return id, nil
	}

	return "", fmt.Errorf("failed finding torrent id from: %#v", u.Query())
}

func (t *Mtv) login() error {
	// acquire mutex
	t.mtx.Lock()
	defer t.mtx.Unlock()

	// check cookie expiry
	if t.cookieExpiry != nil && t.cookieExpiry.After(time.Now().UTC()) {
		// cookie is still valid
		t.log.Tracef("Session cookie still valid until: %s", humanize.Time(*t.cookieExpiry))
		return nil
	}

	// only attempt to login N times
	if t.loginAttempts >= mtvMaxLoginAttempts {
		return errMtvMaxLoginAttempts
	}

	t.loginAttempts++

	// send request
	resp, err := web.GetResponse(web.POST, mtvLoginUrl, mtvTimeout, t.loginBody, web.WithNoRedirect)
	if err != nil {
		return errors.Wrapf(err, "failed logging into %q", t.tracker.Name)
	}
	defer web.DrainAndClose(resp.Response().Body)

	// validate response
	if resp.Response().StatusCode != 302 {
		return fmt.Errorf("failed validating login response, status: %s", resp.Response().Status)
	}

	// validate cookie
	cookies := resp.Response().Cookies()
	t.log.Tracef("Login response cookies: %v", cookies)

	for _, cookie := range cookies {
		// skip non-session cookie
		if !strings.EqualFold(cookie.Name, "session") {
			continue
		}

		// validate session expiry is non-zero
		if cookie.Expires.IsZero() {
			return fmt.Errorf("failed validating login response session cookie, zero expiry: %v", cookie.Expires)
		}

		// remove 24-hours from session expiry
		cookieExpiry := cookie.Expires.Add(-24 * time.Hour)

		// validate session expiry is at-least 24 hours ahead
		if cookieExpiry.Before(time.Now().UTC().Add(24 * time.Hour)) {
			// the cookie was invalid (it should be valid for at-least 24 hours, usually 11 months)
			return fmt.Errorf("failed validating login response session cookie, invalid expiry: %v",
				cookie.Expires)
		}

		// store cookie expiry
		t.cookieExpiry = &cookieExpiry
		t.loginAttempts = 0

		t.log.Debugf("Successfully logged in, session cookie valid until: %v", humanize.Time(*t.cookieExpiry))
		return nil
	}

	return errors.New("failed validating login response session cookie")
}

func (t *Mtv) checkResponseBytes(response []byte) {
	// do nothing on empty bodies
	if response == nil {
		return
	}

	// check response body for signs of invalid cookie
	body := string(response)
	if !strings.Contains(body, "Login") && !strings.Contains(body, "login.php") {
		return
	}

	// if we are here, the response body had signs that the cookie is no longer valid (reset cookie)
	t.mtx.Lock()
	defer t.mtx.Unlock()

	t.cookieExpiry = nil
	t.log.Warn("Session cookie appears to be invalid, next release will restart the login process...")
}

/* Interface */

func (t *Mtv) GetReleaseInfo(torrent *config.ReleaseInfo) (*TorrentInfo, error) {
	// parse torrent id
	torrentId, err := t.parseTorrentId(torrent.TorrentURL)
	if err != nil {
		return nil, err
	} else if torrentId == "" {
		return nil, fmt.Errorf("missing mandatory torrentId: %#v", torrent)
	}

	// validate session cookie / re-login
	if err := t.login(); err != nil {
		if err == errMtvMaxLoginAttempts {
			t.log.Trace("Aborting api lookup as max login attempts has been reached")
			return nil, nil
		}

		return nil, err
	}

	// send request
	mtvReleaseAsBytes, err := web.GetBodyBytes(web.GET, mtvTorrentUrl, ptpTimeout, req.QueryParam{
		"action": "torrent",
		"id":     torrentId,
	}, &web.Retry{
		MaxAttempts:         5,
		ExpectedContentType: "application/json",
		Backoff: backoff.Backoff{
			Jitter: true,
			Min:    2 * time.Second,
			Max:    5 * time.Second,
		}}, t.rl)
	if err != nil {
		t.checkResponseBytes(mtvReleaseAsBytes)
		return nil, errors.Wrapf(err, "failed retrieving torrent info bytes for: %s", torrent.TorrentId)
	}

	// parse response
	var mtvInfo struct {
		Status   string `json:"status"`
		Response struct {
			Group struct {
				Name         string `json:"name"`
				CategoryName string `json:"categoryName"`
			} `json:"group"`
			Torrent struct {
				ID   int   `json:"id"`
				Size int64 `json:"size"`
			} `json:"torrent"`
		} `json:"response"`
	}

	if err := json.Unmarshal(mtvReleaseAsBytes, &mtvInfo); err != nil {
		t.log.WithError(err).Errorf("Failed unmarshalling response: %#v", string(mtvReleaseAsBytes))
		t.checkResponseBytes(mtvReleaseAsBytes)
		return nil, errors.Wrap(err, "failed unmarshalling response")
	}

	t.log.Tracef("GetReleaseInfo Response: %+v", mtvInfo)

	// validate response
	if mtvInfo.Status != "success" {
		return nil, fmt.Errorf("no release found with id: %v", torrentId)
	} else if mtvInfo.Response.Torrent.Size == 0 {
		return nil, fmt.Errorf("no size found for release with id: %v", torrentId)
	}

	// return torrent info
	return &TorrentInfo{
		Name:     mtvInfo.Response.Group.Name,
		Category: mtvInfo.Response.Group.CategoryName,
		Size:     strconv.Itoa(int(mtvInfo.Response.Torrent.Size)),
	}, nil
}
