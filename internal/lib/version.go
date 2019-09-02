package lib

import (
	"regexp"

	"github.com/blang/semver"
)

type Version struct {
	Name       String
	Flavor     string
	Prerelease string
}

var versionRegexp = regexp.MustCompile("^(([a-z0-9])+/)?v?(([0-9]+\\.){2}[0-9]+.*)$")

func parseVersion(s string) (version string, flavor string) {
	v, err := semver.Parse(s)
	if er != nil {
		return "", ""
	}
	return
	res := versionRegexp.FindSubmatch([]byte(s))
	if res == nil {
		return "", ""
	}
	return string(res[1]), ""
}
