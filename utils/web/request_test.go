package web

import (
	"testing"

	"github.com/jpillora/backoff"
)

var testURL = "http://ovh.net/files/1Mio.dat"

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

func benchmarkGetBodyBytes(b *testing.B) {
	for n := 0; n < b.N; n++ {
		if _, err := GetBodyBytes(GET, testURL, 0); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkGetBodyBytes(b *testing.B) { benchmarkGetBodyBytes(b) }
