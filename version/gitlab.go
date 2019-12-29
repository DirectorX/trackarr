package version

import (
	"fmt"
	"time"

	"gitlab.com/cloudb0x/trackarr/utils/web"

	"github.com/imroc/req"
	"github.com/pkg/errors"
)

type GitlabRelease struct {
	Name            string    `json:"name"`
	Tag             string    `json:"tag_name"`
	CreatedAt       time.Time `json:"created_at"`
	ReleasedAt      time.Time `json:"released_at"`
	UpcomingRelease bool      `json:"upcoming_release"`
}

func gitlabReleases(apiUrl string, privateToken string) ([]GitlabRelease, error) {
	// set headers
	headers := req.Header{}
	if privateToken != "" {
		headers["PRIVATE-TOKEN"] = privateToken
	}

	// retrieve latest releases
	resp, err := web.GetResponse(web.GET, apiUrl, 15, headers)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid gitlab releases response from: %q", apiUrl)
	}

	defer resp.Response().Body.Close()

	// validate response
	if resp.Response().StatusCode != 200 {
		return nil, fmt.Errorf("bad gitlab releases response from: %q, resp = %s",
			apiUrl, resp.Response().Status)
	}

	// decode response
	var releases []GitlabRelease
	if err := resp.ToJSON(&releases); err != nil {
		return nil, errors.Wrap(err, "failed decoding gitlab releases response")
	}

	return releases, nil
}
