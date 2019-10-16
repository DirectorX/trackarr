package release

/* Public */

func (r *TrackerRelease) Process() {
	r.Log.Debugf("Processing release: %s", r.TorrentName)
}
