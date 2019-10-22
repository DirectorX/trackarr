package release

import (
	"strings"
	"time"

	"github.com/l3uddz/trackarr/config"
	"github.com/l3uddz/trackarr/utils/maps"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

/* Public */

func FromMap(t *config.TrackerInstance, log *logrus.Entry, vars map[string]string) (*Release, error) {
	release := &Release{
		Tracker: t,
		Log:     log,
		Info: &config.ReleaseInfo{
			TrackerName: *t.Info.ShortName,
			ReleaseTime: time.Now().Format(time.RFC3339),
		},
	}

	// parse mandatory fields
	if torrentName, err := maps.GetFirstStringMapValue(vars, []string{"torrentName", "$torrentName"},
		false); err != nil {
		release.Log.WithError(err).Error("Failed parsing required field from parse match")
		return nil, errors.Wrap(err, "failed parsing required field")
	} else {
		release.Info.TorrentName = torrentName
	}

	if torrentURL, err := maps.GetFirstStringMapValue(vars, []string{"torrentUrl", "$torrentUrl"},
		false); err != nil {
		release.Log.WithError(err).Error("Failed parsing required field from parse match")
		return nil, errors.Wrap(err, "failed parsing required field")
	} else {
		release.Info.TorrentURL = torrentURL
	}

	// parse non-mandatory fields
	if torrentSize, err := maps.GetFirstStringMapValue(vars, []string{"torrentSize", "$torrentSize", "size", "$size"},
		false); err == nil {
		release.Info.SizeString = strings.Replace(torrentSize, ",", "", -1)
	}

	if torrentCategory, err := maps.GetFirstStringMapValue(vars, []string{"$category", "category"},
		false); err == nil {
		release.Info.Category = torrentCategory
	}

	if torrentEncoder, err := maps.GetFirstStringMapValue(vars, []string{"encoder", "$encoder"},
		false); err == nil {
		release.Info.Encoder = torrentEncoder
	}

	if torrentResolution, err := maps.GetFirstStringMapValue(vars, []string{"resolution", "$resolution"},
		false); err == nil {
		release.Info.Resolution = torrentResolution
	}

	if torrentContainer, err := maps.GetFirstStringMapValue(vars, []string{"container", "$container"},
		false); err == nil {
		release.Info.Container = torrentContainer
	}

	if torrentOrigin, err := maps.GetFirstStringMapValue(vars, []string{"origin", "$origin"},
		false); err == nil {
		release.Info.Origin = torrentOrigin
	}

	if torrentTags, err := maps.GetFirstStringMapValue(vars, []string{"$releaseTags", "$tags", "releaseTags", "tags"},
		false); err == nil {
		release.Info.Tags = torrentTags
	}

	return release, nil
}