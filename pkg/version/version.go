package version

import (
	"bytes"
	"fmt"
)

var (
	// GitCommit is the Full Git Commit SHA
	Revision string
	// GitCommitAuthor is the author of the Git Commit
	Author string
	// GitBranch is the Full Git Branch Name
	Branch string
	// BuildDate is the DateTimeStamp during build
	BuildDate string
	// GitDescribe is a way to intentionally describe the version
	GitDescribe string
	// Version is the Full Semantic Version
	Version string
	// VersionPrerelease is the pre-release name (dev,rc-1,alpha,beta,nightly,etc.)
	Prerelease string
	// VersionMetaData is the optional metadata to attach to a version
	Metadata string
	// Builder is the name of the user that builds the artifact (i.e whoami)
	Builder string
	// BuildHost is the name of the host that builds the artifact
	BuildHost string
)

func init() {
	if GitDescribe != "" {
		Version = GitDescribe
	}

	if GitDescribe == "" && Prerelease == "" {
		Prerelease = "dev"
	}
}

func Number() string {
	if Version == "" && Prerelease == "" {
		return "(version unknown)"
	}

	version := fmt.Sprintf("%s", Version)

	if Prerelease != "" {
		version = fmt.Sprintf("%s-%s", version, Prerelease)
	}
	if Metadata != "" {
		version = fmt.Sprintf("%s+%s", version, Metadata)
	}
	return version
}

func FullVersionNumber(rev bool) string {
	var versionString bytes.Buffer

	if Version == "" && Prerelease == "" {
		return "(version unknown)"
	}

	fmt.Fprintf(&versionString, "Eve  %s", Version)

	if Prerelease != "" {
		fmt.Fprintf(&versionString, "-%s", Prerelease)
	}

	if Metadata != "" {
		fmt.Fprintf(&versionString, "+%s", Metadata)
	}

	if rev && Revision != "" {
		fmt.Fprintf(&versionString, " (%s)", Revision)
	}

	return versionString.String()
}
