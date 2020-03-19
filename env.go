package beluga

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"

	"github.com/gobuffalo/envy"
)

const (
	envProduction = "production"
	envStaging    = "staging"
	envReview     = "review"
)

var defaultEnv = Environment{
	varDockerContext:    ".",
	varGitDefaultBranch: "master",
}

var envs = []func(Environment){}

func Env() Environment {
	// Start with the defaults
	e := defaultEnv.clone()

	// Load all environment variables
	_ = envy.Load()
	e.Merge(Environment(envy.Map()))

	envReadEnvOverrides(e)

	// CI environments
	for _, read := range envs {
		read(e)
	}

	// Try to fill in blanks with git
	gitEnvRead(e)

	return e
}

type Environment map[string]string

func (e Environment) Get(key, fallback string) string {
	v := e[key]
	if v == "" {
		return fallback
	}
	return v
}

func (e Environment) SortedKeys() []string {
	keys := sort.StringSlice{}
	for key := range e {
		keys = append(keys, key)
	}
	sort.Sort(keys)
	return []string(keys)
}

func (e Environment) KnownKeys() []string {
	keys := sort.StringSlice{}
	for _, key := range knownVarNames {
		if v, ok := e[key]; ok && v != "" {
			keys = append(keys, key)
		}
	}
	sort.Sort(keys)
	return []string(keys)
}

func (e Environment) Merge(src Environment) {
	for key, value := range src {
		if value != "" {
			e[key] = value
		}
	}
}

func (e Environment) MergeMissing(src Environment) {
	for key, value := range src {
		if e[key] == "" && value != "" {
			e[key] = value
		}
	}
}

func (e Environment) clone() Environment {
	n := make(Environment)
	for key, value := range e {
		if value != "" {
			n[key] = value
		}
	}
	return n
}

// envReadEnvOverrides overrides values for a specific environment by reading
// the value of the BELUGA_* variable as an env file.
func envReadEnvOverrides(e Environment) {
	v := e["BELUGA_"+strings.ToUpper(e[varEnvironment])]
	s, err := godotenv.Unmarshal(v)
	if err != nil {
		log.Println("beluga: env parse: ", err)
	}

	e.Merge(Environment(s))
}

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
