package beluga

import (
	"log"
	"os"
	"sort"
	"strings"

	"github.com/joho/godotenv"

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
	envy.Load()
	e.Merge(Environment(envy.Map()))

	envReadEnvOverrides(e)

	// CI environments
	for _, read := range envs {
		read(e)
	}

	// Try to fill in blanks with git
	gitEnvRead(e)

	// Compute docker-compose.yaml|yml
	if e[varDockerComposeFile] == "" {
		chkFile := func(f string) bool {
			if _, err := os.Stat(f); err == nil {
				e[varDockerComposeFile] = f
				return true
			}
			return false
		}
		switch {
		case chkFile("docker-compose.beluga.yaml"):
		case chkFile("docker-compose.yaml"):
		case chkFile("docker-compose.yml"):
		}
	}

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

func (e Environment) DeployMode() DeployMode {
	return ComposeMode
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
