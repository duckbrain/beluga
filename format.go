package beluga

import (
	"errors"
	"fmt"
	"strings"
)

type EnvFormat int

const (
	EnvFileFormat EnvFormat = iota
	GoEnvFormat
	BashFormat
)

// Format returns the env variables in the specified format. If allKeys is true,
// it will return all values found in the environment, otherwise it only outputs
// keys known to Beluga. Errors may occur for individual key outputs, in these
// cases the function will return a valid slice of strings, and the error.
func (env Environment) Format(format EnvFormat, allKeys bool) ([]string, error) {
	var formatter func(key, val string) (string, error)
	switch format {
	case BashFormat:
		formatter = func(key, val string) (string, error) {
			// TODO escape special characters like newline and single quote
			return fmt.Sprintf("export %v='%v'", key, env[key]), nil
		}
	case EnvFileFormat, GoEnvFormat:
		formatter = func(key, val string) (s string, err error) {
			// ENV file format doesn't support newlines.
			if format == EnvFileFormat && strings.Contains(val, "\n") {
				err = errors.New("cannot represent newlines in env file format")
			}
			s = fmt.Sprintf("%v=%v", key, env[key])
			return
		}
	default:
		return nil, errors.New("unknown format")
	}

	var err error
	var values []string
	if allKeys {
		values = env.SortedKeys()
	} else {
		values = env.KnownKeys()
	}
	for i, key := range values {
		values[i], err = formatter(key, env[key])
	}
	return values, err
}
