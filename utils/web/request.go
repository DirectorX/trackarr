package web

import (
	"github.com/imroc/req"
	"github.com/pkg/errors"
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

func GetResponse(method HTTPMethod, requestUrl string, timeout int, v ...interface{}) (*req.Resp, error) {
	var resp *req.Resp
	var err error

	// prepare request
	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}

	inputs := make([]interface{}, 0)
	inputs = append(inputs, client)
	inputs = append(inputs, v...)

	// send request
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
		return nil, err
	} else {
		log.Tracef("Request URL: %s", resp.Request().URL)
		log.Tracef("Request Response: %s", resp.Response().Status)
	}

	return resp, nil
}

func GetBodyBytes(method HTTPMethod, requestUrl string, timeout int, v ...interface{}) ([]byte, error) {
	// send request
	resp, err := GetResponse(method, requestUrl, timeout, v...)
	if err != nil {
		return nil, err
	}

	// process response
	bodyBytes, err := resp.ToBytes()
	if err != nil {
		log.WithError(err).Errorf("Failed reading response body for url: %q", requestUrl)
		return nil, errors.Wrap(err, "failed reading url response body")
	}

	return bodyBytes, nil
}

func GetBodyString(method HTTPMethod, requestUrl string, timeout int, v ...interface{}) (string, error) {
	bodyBytes, err := GetBodyBytes(method, requestUrl, timeout, v...)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}
