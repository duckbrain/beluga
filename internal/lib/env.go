package lib

import (
	"os"
	"sort"
	"strings"
)

type EnvReader interface {
	EnvRead() Env
}

type Env map[string]string

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

func (e Env) Merge(source Env) {
	if source == nil {
		return
	}
	for key, value := range source {
		if e[key] == "" && value != "" {
			e[key] = value
		}
	}
}

var envs = compositeEnv{belugaEnv("")}

func ParseEnv() Env {
	return envs.EnvRead()
}

type compositeEnv []EnvReader

func (envs compositeEnv) EnvRead() Env {
	vals := make(Env)
	for _, envReader := range envs {
		env := envReader.EnvRead()
		vals.Merge(env)
	}
	return vals
}

type belugaEnv string

func (prefix belugaEnv) EnvRead() Env {
	vals := make(Env)
	for _, line := range os.Environ() {
		i := strings.IndexRune(line, '=')
		key := line[:i]
		value := line[i+1:]
		if strings.HasPrefix(key, string(prefix)) {
			vals[key] = value
		}
	}
	return vals
}
