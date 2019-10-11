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

// GetBody - Retrieve the body of a web page as a string
func GetBody(method HTTPMethod, url string, timeout int) (string, error) {
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
		return "", errors.New("request method has not been implemented")
	}

	log.Tracef("Request URL: %s", resp.Request.URL)
	log.Tracef("Request Response: %s", resp.Status)

	// validate response
	if err != nil {
		log.WithError(err).Errorf("Failed retrieving body for page: %q", url)
		return "", errors.Wrap(err, "failed retrieving page body")
	}

	// process response
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Errorf("Failed reading response body for page: %q", url)
		return "", errors.Wrap(err, "failed reading page response body")
	}

	body := string(bodyBytes)
	return body, nil
}
