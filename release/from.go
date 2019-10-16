package release

import (
	"github.com/l3uddz/trackarr/autodl/processor"
	"github.com/l3uddz/trackarr/utils/maps"
	"github.com/pkg/errors"
)

/* Structs */

type TrackerRelease struct {
	TrackerName       string
	TorrentName       *string
	TorrentURL        *string
	TorrentSizeString *string
	TorrentCategory   *string
}

/* Public */

func FromMap(p *processor.Processor, vars *map[string]string) (*TrackerRelease, error) {
	release := &TrackerRelease{TrackerName: p.Tracker.LongName}

	// parse mandatory fields
	if torrentName, err := maps.GetFirstStringMapValue(vars, []string{"torrentName", "$torrentName"},
		false); err != nil {
		p.Log.WithError(err).Errorf("Failed parsing required field %q from parse match", "torrentName")
		return nil, errors.Wrap(err, "failed parsing required field torrentName")
	} else {
		release.TorrentName = &torrentName
	}

	if torrentURL, err := maps.GetFirstStringMapValue(vars, []string{"torrentUrl", "$torrentUrl"},
		false); err != nil {
		p.Log.WithError(err).Errorf("Failed parsing required field %q from parse match", "torrentUrl")
		return nil, errors.Wrap(err, "failed parsing required field torrentUrl")
	} else {
		release.TorrentURL = &torrentURL
	}

	// parse non-mandatory fields
	if torrentSize, err := maps.GetFirstStringMapValue(vars, []string{"torrentSize", "$torrentSize", "size", "$size"},
		false); err != nil {
		p.Log.WithError(err).Tracef("Failed parsing field %q from parse match", "torrentSize")
	} else {
		release.TorrentSizeString = &torrentSize
	}

	if torrentCategory, err := maps.GetFirstStringMapValue(vars, []string{"$category", "category", "torrentCategory",
		"$torrentCategory"}, false); err != nil {
		p.Log.WithError(err).Tracef("Failed parsing field %q from parse match", "category")
	} else {
		release.TorrentCategory = &torrentCategory
	}

	return release, nil
}
