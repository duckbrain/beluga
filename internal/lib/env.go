package lib

import (
	"log"
	"path/filepath"
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

var Environment = Env{
	varDockerContext:    ".",
	varGitDefaultBranch: "master",
}

func init() {
	// Start with the defaults
	e := Environment
	// BELUGA_ prefix vars
	belugaEnv{}.EnvRead(e)
	// Environment specific overrides
	belugaEnvironmentEnv(e.Environment()).EnvRead(e)

	// CI environments
	envs.EnvRead(e)

	// Compute Dockerfile if not yet set
	if e[varDockerfile] == "" {
		e[varDockerfile] = filepath.Join(e.DockerContext(), "Dockerfile")
	}
}

type EnvReader interface {
	EnvRead(Env)
}

type Env map[string]string

func (e Env) Get(key, fallback string) string {
	v := e[key]
	if v == "" {
		return fallback
	}
	return v
}

func (e Env) SortedKeys() []string {
	keys := sort.StringSlice{}
	for key := range e {
		keys = append(keys, key)
	}
	sort.Sort(keys)
	return []string(keys)
}

func (e Env) KnownKeys() []string {
	keys := sort.StringSlice{}
	for _, key := range knownVarNames {
		if v, ok := e[key]; ok && v != "" {
			keys = append(keys, key)
		}
	}
	sort.Sort(keys)
	return []string(keys)
}

func (e Env) Merge(src Env) {
	for key, value := range src {
		if value != "" {
			e[key] = value
		}
	}
}

func (e Env) MergeMissing(src Env) {
	for key, value := range src {
		if e[key] == "" && value != "" {
			e[key] = value
		}
	}
}

var envs = compositeEnv{}

type compositeEnv []EnvReader

func (c compositeEnv) EnvRead(e Env) {
	for _, r := range c {
		r.EnvRead(e)
	}
}

type belugaEnv struct{}

func (belugaEnv) EnvRead(e Env) {
	envy.Load()
	e.Merge(Env(envy.Map()))
}

type belugaEnvironmentEnv string

func (env belugaEnvironmentEnv) EnvRead(e Env) {
	v := e["BELUGA_"+strings.ToUpper(e.Environment())]
	s, err := godotenv.Unmarshal(v)
	if err != nil {
		log.Println("beluga: env parse: ", err)
	}

	e.Merge(Env(s))
}
