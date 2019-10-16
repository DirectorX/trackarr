package release

import "github.com/l3uddz/trackarr/autodl/processor"

/* Public */

func Process(p *processor.Processor, release *TrackerRelease) {
	p.Log.Debugf("Processing release: %s", release.TorrentName)
}
