package main

import (
	"fmt"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

// Version is needed for ggallin
const Version = "1.0.0"

// GitCommit is needed for ggallin
var GitCommit = ""

func doSelfupdate() error {
	v := semver.MustParse(Version)
	latest, err := selfupdate.UpdateSelf(v, "delphinus/git-commit-autouser")
	if err != nil {
		return fmt.Errorf("binary update failed: %v", err)
	}
	if latest.Version.Equals(v) {
		fmt.Println("current binary is the latest version", Version)
	} else {
		fmt.Println("successfully updated to version", latest.Version)
		fmt.Println("release note:\n", latest.ReleaseNotes)
	}
	return nil
}
