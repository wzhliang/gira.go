package release

import (
	"fmt"

	"github.com/blang/semver/v4"
)

// Release ...
type Release struct {
	raw     string
	Major   uint64
	Minor   uint64
	Patch   uint64
	Project string
}

// New ...
func New(fixVersion string) *Release {
	sv, err := semver.Make(fixVersion[1:])
	if err != nil {
		return nil
	}
	return &Release{
		raw:     fixVersion,
		Major:   sv.Major,
		Minor:   sv.Minor,
		Patch:   sv.Patch,
		Project: "", // TODO: add project support
	}
}

func (r *Release) shouldGotoMaster() bool {
	return r.Patch == 0
}

// TargetBranch ...
func (r *Release) TargetBranch() string {
	if r.shouldGotoMaster() {
		return "master"
	}
	// TODO: add project support
	return fmt.Sprintf("release-%d.%d", r.Major, r.Minor)
}
