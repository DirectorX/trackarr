package web

import (
	"net/http"
)

type Option int

const (
	WithNoRedirect Option = iota + 1
)

/* Private */

func setOption(opt Option, client *http.Client) {
	switch opt {
	case WithNoRedirect:
		noRedirect(client)
	default:
		break
	}
}

/* Public */

func noRedirect(client *http.Client) {
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
}
