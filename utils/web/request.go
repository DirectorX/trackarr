package web

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/l3uddz/trackarr/logger"
)

var (
	log = logger.GetLogger("web")
)

/* Structs */

// HTTPMethod - The HTTP request method to use
type HTTPMethod int

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

func GetBodyBytes(method HTTPMethod, url string, timeout int) ([]byte, error) {
	var resp *http.Response
	var err error

	// create client
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	// send request
	switch method {
	case GET:
		resp, err = client.Get(url)
	default:
		log.Error("Request method has not been implemented")
		return nil, errors.New("request method has not been implemented")
	}

	// validate response
	if err != nil {
		log.WithError(err).Errorf("Failed retrieving body for page: %q", url)
		return nil, errors.Wrap(err, "failed retrieving page body")
	} else {
		log.Tracef("Request URL: %s", resp.Request.URL)
		log.Tracef("Request Response: %s", resp.Status)
	}

	// process response
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Errorf("Failed reading response body for page: %q", url)
		return nil, errors.Wrap(err, "failed reading page response body")
	}

	return bodyBytes, nil
}

func GetBodyString(method HTTPMethod, url string, timeout int) (string, error) {
	bodyBytes, err := GetBodyBytes(method, url, timeout)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}
