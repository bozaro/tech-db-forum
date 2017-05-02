package tests

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
)

var (
	Project   = "tech-db-forum"
	Version   = "0.3.0"
	BuildTag  string
	GitCommit string
)

const (
	VERSION_UNKNOWN = iota
	VERSION_OUTDATE
	VERSION_LOCAL
	VERSION_LATEST
)

type githubRef struct {
	Ref    string
	Url    string
	Object githubObject
}

type githubObject struct {
	Sha  string
	Type string
	Url  string
}

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

func VersionCheck() (int, error) {
	if GitCommit == "" {
		return VERSION_LOCAL, nil
	}

	req, err := http.NewRequest("GET", "https://api.github.com/repos/bozaro/tech-db-forum/git/refs/heads/master", nil)
	if err != nil {
		return VERSION_UNKNOWN, err
	}

	res, err := HttpTransport.RoundTrip(req)
	if err != nil {
		return VERSION_UNKNOWN, err
	} else if res.StatusCode != 200 {
		return VERSION_UNKNOWN, errors.New("Unexpected request status code: " + res.Status)
	}

	payload, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return VERSION_UNKNOWN, err
	}

	ref := githubRef{}
	err = json.Unmarshal(payload, &ref)
	if err != nil {
		return VERSION_UNKNOWN, err
	}

	if ref.Object.Sha == GitCommit {
		return VERSION_LATEST, nil
	}
	return VERSION_OUTDATE, nil
}
