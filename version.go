package beluga

import (
	"regexp"
)

var versionRegexp = regexp.MustCompile(
	`^((\w+)\/)?v?((\d+\.\d+\.\d+)(-(\w+)(\.[0-9])?)?)$`,
)

func parseVersion(s string) Environment {
	res := versionRegexp.FindStringSubmatch(s)
	if res == nil {
		return Environment{}
	}
	// Uncomment to debug the array that's output in tests
	// x, _ := json.Marshal(res)
	// fmt.Println(string(x))
	return Environment{
		varApplication: res[2],
		varVersion:     res[3],
		varEnvironment: res[6],
	}
}
