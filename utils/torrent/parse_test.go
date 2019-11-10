package torrent

import (
	"testing"

	"github.com/l3uddz/trackarr/cache"

	"github.com/imroc/req"
)

func TestGetTorrentDetails(t *testing.T) {
	type args struct {
		torrentUrl string
		timeout    int
		headers    req.Header
	}
	tests := []struct {
		name    string
		args    args
		want    *Data
		wantErr bool
	}{
		{
			name: "basic",
			args: args{
				torrentUrl: "http://releases.ubuntu.com/19.10/ubuntu-19.10-live-server-amd64.iso.torrent",
			},
			wantErr: false,
		},
	}

	if err := cache.Init(); err != nil {
		t.Errorf("Failed to init cache: %s", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTorrentDetails(tt.args.torrentUrl, tt.args.timeout, tt.args.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTorrentDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("GetTorrentDetails() = %v, want %v", got, tt.want)
			// }
			t.Log(got)
		})
	}
}
