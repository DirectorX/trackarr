package version

import (
	"github.com/Masterminds/semver"
	"github.com/l3uddz/trackarr/logger"
)

var (
	log = logger.GetLogger("ver")
)

func IsLatestGitlabVersion(apiUrl string, privateToken string, currentVersion string) (bool, string) {
	// create semver from currentVersion
	cVer, err := semver.NewVersion(currentVersion)
	if err != nil {
		log.WithError(err).Errorf("Failed creating semver from currentVersion: %q", currentVersion)
		return true, currentVersion
	}

	// retrieve latest releases
	releases, err := gitlabReleases(apiUrl, privateToken)
	if err != nil {
		log.WithError(err).Error("Failed retrieving latest Gitlab releases...")
		return true, currentVersion
	}
	log.Debugf("Found %d Gitlab releases", len(releases))

	// iterate releases
	highestVersion := &semver.Version{}

	for _, release := range releases {
		// translate release tag to semver
		v, err := semver.NewVersion(release.Tag)
		if err != nil {
			log.WithError(err).Errorf("Failed converting release tag to version: %q", release.Tag)
			continue
		}

		// are we interested in this version?
		if v.Prerelease() != cVer.Prerelease() {
			// skip this release
			continue
		}

		// is this greater than our current version?
		if v.GreaterThan(cVer) && v.GreaterThan(highestVersion) {
			highestVersion = v
		}
	}

	if highestVersion != nil {
		return false, highestVersion.String()
	}

	return true, currentVersion
}
