package version

import (
	"testing"

	"github.com/Masterminds/semver"
)

func TestVersion_BranchGreater(t *testing.T) {
	type fields struct {
		Current      *semver.Version
		apiUrl       string
		privateToken string
	}

	tests := []struct {
		name        string
		version     *semver.Version
		rVersion    *semver.Version
		wantGreater bool
		wantBranch  bool
	}{
		{
			name:        "sameDev",
			version:     semver.MustParse("v0.1.0-dev"),
			rVersion:    semver.MustParse("v0.1.0-dev"),
			wantGreater: false,
			wantBranch:  true,
		},
		{
			name:        "sameMaster",
			version:     semver.MustParse("v0.1.0"),
			rVersion:    semver.MustParse("v0.1.0"),
			wantGreater: false,
			wantBranch:  true,
		},
		{
			name:        "lowerDev",
			version:     semver.MustParse("v1.0.0-dev"),
			rVersion:    semver.MustParse("v0.2.0-dev"),
			wantGreater: false,
			wantBranch:  true,
		},
		{
			name:        "higherDev",
			version:     semver.MustParse("v0.1.0-dev"),
			rVersion:    semver.MustParse("v1.0.0-dev"),
			wantGreater: true,
			wantBranch:  true,
		},
		{
			name:        "lowerMaster",
			version:     semver.MustParse("v1.0.0"),
			rVersion:    semver.MustParse("v0.2.0"),
			wantGreater: false,
			wantBranch:  true,
		},
		{
			name:        "higherMaster",
			version:     semver.MustParse("v0.1.0"),
			rVersion:    semver.MustParse("v1.0.0"),
			wantGreater: true,
			wantBranch:  true,
		},
		{
			name:        "lowerDevBranch",
			version:     semver.MustParse("v1.1.0-dev"),
			rVersion:    semver.MustParse("v1.0.0"),
			wantGreater: false,
			wantBranch:  false,
		},
		{
			name:        "higherDevBranch",
			version:     semver.MustParse("v0.1.0-dev"),
			rVersion:    semver.MustParse("v1.0.0"),
			wantGreater: true,
			wantBranch:  false,
		},
		{
			name:        "lowerMasterBranch",
			version:     semver.MustParse("v1.1.0"),
			rVersion:    semver.MustParse("v1.0.0-dev"),
			wantGreater: false,
			wantBranch:  false,
		},
		{
			name:        "higherMasterBranch",
			version:     semver.MustParse("v0.1.0"),
			rVersion:    semver.MustParse("v1.0.0-dev"),
			wantGreater: true,
			wantBranch:  false,
		},
		{
			name:        "lowerDevCommit",
			version:     semver.MustParse("v1.1.0-dev+1234"),
			rVersion:    semver.MustParse("v1.0.0-dev"),
			wantGreater: false,
			wantBranch:  true,
		},
		{
			name:        "higherDevCommit",
			version:     semver.MustParse("v1.1.0-dev+1234"),
			rVersion:    semver.MustParse("v1.2.0-dev"),
			wantGreater: true,
			wantBranch:  true,
		},
		{
			name:        "higherDevCommitRemote",
			version:     semver.MustParse("v1.1.0-dev+1234"),
			rVersion:    semver.MustParse("v1.2.0-dev+abcd"),
			wantGreater: true,
			wantBranch:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Version{
				Current: tt.version,
				apiUrl:  apiUrl,
			}
			greater := v.isGreater(tt.rVersion)
			if greater != tt.wantGreater {
				t.Errorf("Version.isGreater() got = %v, want %v", greater, tt.wantGreater)
			}
			sameBranch := v.isSameBranch(tt.rVersion)
			if sameBranch != tt.wantBranch {
				t.Errorf("Version.isSameBranch() got = %v, want %v", sameBranch, tt.wantBranch)
			}
		})
	}
}
