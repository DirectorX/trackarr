package release

import (
	"github.com/l3uddz/trackarr/autodl/parser"
	"github.com/l3uddz/trackarr/config"
	"github.com/l3uddz/trackarr/utils/maps"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

/* Structs */

type TrackerRelease struct {
	Tracker           *parser.TrackerInfo
	Log               *logrus.Entry
	Cfg               *config.TrackerConfiguration
	TorrentName       string
	TorrentURL        string
	TorrentSizeString *string
	TorrentSizeBytes  *int64
	TorrentCategory   *string
}

/* Public */

func FromMap(t *parser.TrackerInfo, c *config.TrackerConfiguration,
	log *logrus.Entry, vars *map[string]string) (*TrackerRelease, error) {
	release := &TrackerRelease{Tracker: t, Cfg: c, Log: log}

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
		release.TorrentSizeString = &torrentSize
	}

	if torrentCategory, err := maps.GetFirstStringMapValue(vars, []string{"$category", "category", "torrentCategory",
		"$torrentCategory"}, false); err == nil {
		release.TorrentCategory = &torrentCategory
	}

	return release, nil
}
