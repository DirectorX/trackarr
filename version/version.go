package version

import (
	"gitlab.com/cloudb0x/trackarr/config"
	"gitlab.com/cloudb0x/trackarr/logger"

	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"
)

var (
	log      = logger.GetLogger("ver")
	Trackarr *Version
)

const (
	apiUrl = "https://gitlab.com/api/v4/projects/15385789/releases"
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
		apiUrl:  apiUrl,
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
	for _, release := range releases {
		// translate release tag to semver
		rV, err := semver.NewVersion(release.Tag)
		if err != nil {
			log.WithError(err).Errorf("Failed converting release tag to version: %s", release.Tag)
			continue
		}

		if v.isSameBranch(rV) {
			if v.isGreater(rV) {
				return false, rV
			}

			// Most recent release in the same branch isn't newer, discard everything else
			break
		}
	}

	return true, v.Current
}

func (v *Version) isSameBranch(rV *semver.Version) bool {
	return v.Current.Prerelease() == rV.Prerelease()
}

func (v *Version) isGreater(rV *semver.Version) bool {
	// is this greater than our current version?
	return rV.GreaterThan(v.Current)
}
