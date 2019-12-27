package version

import (
	"github.com/Masterminds/semver"
	"github.com/l3uddz/trackarr/config"
	"github.com/l3uddz/trackarr/logger"
	"github.com/pkg/errors"
)

var (
	log      = logger.GetLogger("ver")
	Trackarr *Version
)

type Version struct {
	Current      *semver.Version
	apiUrl       string
	privateToken string
}

// init sets the exported object
func Init(buildConfig *config.BuildVars) error {
	// Parse current version
	c, err := semver.NewVersion(buildConfig.Version)
	if err != nil {
		return errors.Wrapf(err, "Failed creating semver from currentVersion: %s", buildConfig.Version)
	}

	Trackarr = &Version{
		apiUrl:  "https://gitlab.com/api/v4/projects/15385789/releases",
		Current: c,
	}

	return nil
}

func (v *Version) IsLatest() (bool, *semver.Version) {
	// retrieve latest releases
	releases, err := gitlabReleases(v.apiUrl, v.privateToken)
	if err != nil {
		log.WithError(err).Error("Failed retrieving latest Gitlab releases...")
		return true, v.Current
	}
	log.Debugf("Found %d Gitlab releases", len(releases))

	// iterate releases
	highestVersion := &semver.Version{}

	for _, release := range releases {
		// translate release tag to semver
		rV, err := semver.NewVersion(release.Tag)
		if err != nil {
			log.WithError(err).Errorf("Failed converting release tag to version: %q", release.Tag)
			continue
		}

		// are we interested in this version?
		// skip non matching pre releases
		if (v.Current.Prerelease() == "" && rV.Prerelease() != "") ||
			(v.Current.Prerelease() != "" && rV.Prerelease() == "") {
			continue
		}

		// is this greater than our current version?
		if rV.GreaterThan(v.Current) && rV.GreaterThan(highestVersion) {
			highestVersion = rV
		}
	}

	if highestVersion != nil && highestVersion.String() != "0.0.0" {
		return false, highestVersion
	}

	return true, v.Current
}
