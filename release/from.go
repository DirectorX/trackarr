package release

import (
	"github.com/l3uddz/trackarr/autodl/parser"
	"github.com/l3uddz/trackarr/config"
	"github.com/l3uddz/trackarr/utils/maps"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

/* Structs */

type TrackerRelease struct {
	Tracker     *parser.TrackerInfo
	Log         *logrus.Entry
	Cfg         *config.TrackerConfiguration
	TrackerName string
	ReleaseTime string
	TorrentName string
	TorrentURL  string
	SizeString  string
	SizeBytes   int64
	Category    string
	Encoder     string
	Resolution  string
	Container   string
	Origin      string
	Tags        string
}

/* Public */

func FromMap(t *parser.TrackerInfo, c *config.TrackerConfiguration,
	log *logrus.Entry, vars *map[string]string) (*TrackerRelease, error) {
	release := &TrackerRelease{Tracker: t, Cfg: c, Log: log, TrackerName: *t.ShortName,
		ReleaseTime: time.Now().Format(time.RFC3339)}

	// parse mandatory fields
	if torrentName, err := maps.GetFirstStringMapValue(vars, []string{"torrentName", "$torrentName"},
		false); err != nil {
		release.Log.WithError(err).Error("Failed parsing required field from parse match")
		return nil, errors.Wrap(err, "failed parsing required field")
	} else {
		release.TorrentName = torrentName
	}

	if torrentURL, err := maps.GetFirstStringMapValue(vars, []string{"torrentUrl", "$torrentUrl"},
		false); err != nil {
		release.Log.WithError(err).Error("Failed parsing required field from parse match")
		return nil, errors.Wrap(err, "failed parsing required field")
	} else {
		release.TorrentURL = torrentURL
	}

	// parse non-mandatory fields
	if torrentSize, err := maps.GetFirstStringMapValue(vars, []string{"torrentSize", "$torrentSize", "size", "$size"},
		false); err == nil {
		release.SizeString = strings.Replace(torrentSize, ",", "", -1)
	}

	if torrentCategory, err := maps.GetFirstStringMapValue(vars, []string{"$category", "category"},
		false); err == nil {
		release.Category = torrentCategory
	}

	if torrentEncoder, err := maps.GetFirstStringMapValue(vars, []string{"encoder", "$encoder"},
		false); err == nil {
		release.Encoder = torrentEncoder
	}

	if torrentResolution, err := maps.GetFirstStringMapValue(vars, []string{"resolution", "$resolution"},
		false); err == nil {
		release.Resolution = torrentResolution
	}

	if torrentContainer, err := maps.GetFirstStringMapValue(vars, []string{"container", "$container"},
		false); err == nil {
		release.Container = torrentContainer
	}

	if torrentOrigin, err := maps.GetFirstStringMapValue(vars, []string{"origin", "$origin"},
		false); err == nil {
		release.Origin = torrentOrigin
	}

	if torrentTags, err := maps.GetFirstStringMapValue(vars, []string{"$releaseTags", "$tags", "releaseTags", "tags"},
		false); err == nil {
		release.Tags = torrentTags
	}

	return release, nil
}
