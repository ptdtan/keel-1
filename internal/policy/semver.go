package policy

import (
	"fmt"

	"github.com/Masterminds/semver"
)

// SemverPolicyType - policy type
type SemverPolicyType int

// available policies
const (
	SemverPolicyTypeNone SemverPolicyType = iota
	SemverPolicyTypeAll
	SemverPolicyTypeMajor
	SemverPolicyTypeMinor
	SemverPolicyTypePatch
	// SemverPolicyTypeForce // update always when a new image is available
)

func (t SemverPolicyType) String() string {
	switch t {
	case SemverPolicyTypeNone:
		return "none"
	case SemverPolicyTypeAll:
		return "all"
	case SemverPolicyTypeMajor:
		return "major"
	case SemverPolicyTypeMinor:
		return "minor"
	case SemverPolicyTypePatch:
		return "patch"
	// case PolicyTypeForce:
	// 	return "force"
	default:
		return ""
	}
}

func NewSemverPolicy(spt SemverPolicyType) *SemverPolicy {
	return &SemverPolicy{
		spt: spt,
	}
}

type SemverPolicy struct {
	spt SemverPolicyType
}

func (sp *SemverPolicy) ShouldUpdate(current, new string) (bool, error) {

	if current == "latest" {
		return true, nil
	}

	currentVersion, err := semver.NewVersion(current)
	if err != nil {
		return false, fmt.Errorf("failed to parse current version: %s", err)
	}

	newVersion, err := semver.NewVersion(new)
	if err != nil {
		return false, fmt.Errorf("failed to parse new version: %s", err)
	}

	if currentVersion.Prerelease() != newVersion.Prerelease() && sp.spt != SemverPolicyTypeAll {
		return false, nil
	}

	// new version is not higher than current - do nothing
	if !currentVersion.LessThan(newVersion) {
		return false, nil
	}

	switch sp.spt {
	case SemverPolicyTypeAll, SemverPolicyTypeMajor:
		return true, nil
	case SemverPolicyTypeMinor:
		return newVersion.Major() == currentVersion.Major(), nil
	case SemverPolicyTypePatch:
		return newVersion.Major() == currentVersion.Major() && newVersion.Minor() == currentVersion.Minor(), nil
	}
	return false, nil
}

func (sp *SemverPolicy) Name() string {
	return sp.spt.String()
}

func (sp *SemverPolicy) Type() PolicyType { return PolicyTypeSemver }
