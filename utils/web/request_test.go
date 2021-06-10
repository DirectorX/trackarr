package web

import (
	"strings"
	"testing"

	"github.com/jpillora/backoff"
)

var testURL = "http://proof.ovh.net/files/1Mio.dat"

func TestGetResponse(t *testing.T) {
	type args struct {
		method     HTTPMethod
		requestUrl string
		timeout    int
		v          []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "basicGET",
			args: args{
				method:     GET,
				requestUrl: testURL,
				timeout:    0,
			},
			wantErr: false,
		},
		{
			name: "retryGET",
			args: args{
				method:     GET,
				requestUrl: testURL,
				timeout:    0,
				v: []interface{}{
					&Retry{
						MaxAttempts: 2,
						Backoff: backoff.Backoff{
							Jitter: true,
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetResponse(tt.args.method, tt.args.requestUrl, tt.args.timeout, tt.args.v...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got.Response().StatusCode != 200 {
				t.Errorf("GetResponse() StatusCode = %v, wanted %v", got.Response().StatusCode, 200)
			}
		})
	}
}

func TestCookieLeak(t *testing.T) {
	// set cookies in first request
	body, err := GetBodyString(GET, "https://httpbin.org/cookies/set/test/leaked", 5)
	if err != nil {
		t.Errorf("GetBodyString() error = %v", err)
		return
	} else if !strings.Contains(body, "leaked") {
		t.Errorf("GetBodyString() unexpected set-cookie response = %v", body)
		return
	}

	// get cookies from subsequent request (cookies set above should not be there)
	body, err = GetBodyString(GET, "https://postman-echo.com/cookies", 5)
	if err != nil {
		t.Errorf("GetBodyString() error = %v", err)
		return
	} else if strings.Contains(body, "leaked") {
		t.Errorf("GetBodyString() cookies leaked from initial request = %v", body)
		return
	}

	// get cookies from first request (cookies set initially should be there)
	body, err = GetBodyString(GET, "https://httpbin.org/cookies", 5)
	if err != nil {
		t.Errorf("GetBodyString() error = %v", err)
		return
	} else if !strings.Contains(body, "leaked") {
		t.Errorf("GetBodyString() cookies not retained from initial set request = %v", body)
		return
	}
}

//func TestNoRedirect(t *testing.T) {
//	// test noredirect
//	body, err := GetBodyString(GET, "https://httpbin.org/redirect-to?url=https%3A%2F%2Fwww.google.com", 5,
//		WithNoRedirect)
//	if err != nil {
//		t.Errorf("GetBodyString() error = %v", err)
//		return
//	} else if strings.Contains(body, "Google") {
//		t.Errorf("GetBodyString() unexpected no-redirect response = %v", body)
//		return
//	}
//
//	// test redirect
//	body, err = GetBodyString(GET, "https://httpbin.org/redirect-to?url=https%3A%2F%2Fwww.google.com", 5)
//	if err != nil {
//		t.Errorf("GetBodyString() error = %v", err)
//		return
//	} else if !strings.Contains(body, "Google") {
//		t.Errorf("GetBodyString() unexpected redirect response = %v", body)
//		return
//	}
//}

func benchmarkGetBodyBytes(b *testing.B) {
	for n := 0; n < b.N; n++ {
		if _, err := GetBodyBytes(GET, testURL, 0); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkGetBodyBytes(b *testing.B) { benchmarkGetBodyBytes(b) }
