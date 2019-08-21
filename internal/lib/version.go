package lib

import (
	"regexp"
)

var versionRegexp = regexp.MustCompile("^v?(([0-9]+\\.){2}[0-9]+.*)$")

func parseVersion(s string) string {
	res := versionRegexp.FindSubmatch([]byte(s))
	if res == nil {
		return ""
	}
	return string(res[1])
}
