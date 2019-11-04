package web

import (
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/imroc/req"
	"github.com/jpillora/backoff"
	"github.com/pkg/errors"

	"github.com/l3uddz/trackarr/logger"
	"github.com/l3uddz/trackarr/utils/lists"
)

var (
	// Logging
	log = logger.GetLogger("web")

	// HTTP client
	httpClient = &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 10 * time.Second,
			}).Dial,
		},
	}
)

/* Structs */

// HTTPMethod - The HTTP request method to use
type HTTPMethod int
type Retry struct {
	backoff.Backoff
	MaxAttempts          float64
	RetryableStatusCodes []int
}

const (
	// GET - Use GET HTTP method
	GET HTTPMethod = iota + 1
	// POST - Use POST HTTP method
	POST
	// PUT - Use PUT HTTP method
	PUT
	// DELETE - Use DELETE HTTP method
	DELETE
)

/* Public */

func GetResponse(method HTTPMethod, requestUrl string, timeout int, v ...interface{}) (*req.Resp, error) {
	// prepare request
	client := httpClient
	client.Timeout = time.Duration(timeout) * time.Second

	req.SetJSONEscapeHTML(false)

	inputs := make([]interface{}, 0)
	inputs = append(inputs, client)

	// Extract Retry struct, append everything else
	var retry Retry
	for _, vv := range v {
		switch vT := vv.(type) {
		case *Retry:
			retry = *vT
		case Retry:
			retry = vT
			log.Debugf("Using retry: %#v", retry)
		default:
			inputs = append(inputs, vT)
		}
	}

	// Response var
	var resp *req.Resp
	// Exponential backoff
	for {
		var err error
		switch method {
		case GET:
			resp, err = req.Get(requestUrl, inputs...)
		case POST:
			resp, err = req.Post(requestUrl, inputs...)
		default:
			log.Error("Request method has not been implemented")

			return nil, errors.New("request method has not been implemented")
		}

		// validate response
		if err != nil {
			log.WithError(err).Errorf("Failed requesting url: %q", requestUrl)
			if os.IsTimeout(err) {
				if retry.MaxAttempts == 0 || retry.Attempt() >= retry.MaxAttempts {
					return nil, err
				}

				d := retry.Duration()
				log.Debugf("Retrying failed HTTP request in %s: %q", d, requestUrl)

				time.Sleep(d)
				continue
			}

			return nil, err
		}

		log.Tracef("Request URL: %s", resp.Request().URL)
		log.Tracef("Request Response: %s", resp.Response().Status)

		if retry.MaxAttempts == 0 || retry.Attempt() > retry.MaxAttempts {
			break
		}

		if lists.IntListContains(resp.Response().StatusCode, retry.RetryableStatusCodes) {
			d := retry.Duration()
			log.Debugf("Retrying failed HTTP request in %s: %d - %q", d, resp.Response().StatusCode, requestUrl)

			time.Sleep(d)
			continue
		}

		break
	}

	return resp, nil
}

func GetBodyBytes(method HTTPMethod, requestUrl string, timeout int, v ...interface{}) ([]byte, error) {
	// send request
	resp, err := GetResponse(method, requestUrl, timeout, v...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Response().Body.Close(); err != nil {
			log.WithError(err).Errorf("Failed to close HTTP response body for url: %q", requestUrl)
		}
	}()

	// process response
	body, err := ioutil.ReadAll(resp.Response().Body)
	if err != nil {
		log.WithError(err).Errorf("Failed reading response body for url: %q", requestUrl)
		return nil, errors.Wrap(err, "failed reading url response body")
	}

	return body, nil
}

func GetBodyString(method HTTPMethod, requestUrl string, timeout int, v ...interface{}) (string, error) {
	bodyBytes, err := GetBodyBytes(method, requestUrl, timeout, v...)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}
