package torrent

import (
	"testing"

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
			name:    "basic",
			args:    args{},
			wantErr: false,
		},
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
