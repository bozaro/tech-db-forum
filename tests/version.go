package tests

import (
	"fmt"
	"runtime"
)

var (
	Project   = "tech-db-forum"
	Version   = "0.3.0"
	BuildTag  string
	GitCommit string
)

func VersionFull() string {
	version := fmt.Sprintf("%s/%s", Project, Version)
	if len(GitCommit) >= 7 {
		version += "#" + GitCommit[0:7]
	}
	version += fmt.Sprintf(" (%s %s; %s", runtime.GOOS, runtime.GOARCH, runtime.Version())
	if BuildTag != "" {
		version += fmt.Sprintf("; %s", BuildTag)
	}
	version += ")"
	return version
}
