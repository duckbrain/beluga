package lib

import (
	"os"
	"strings"
)

type Env map[string]string

type EnvReader interface {
	EnvRead() Env
}

type compositeEnv []EnvReader

func (envs compositeEnv) EnvRead() Env {
	vals := make(Env)
	for _, envReader := range envs {
		env := envReader.EnvRead()
		if env == nil {
			continue
		}
		for key, val := range env {
			currentVal := vals[key]
			if currentVal == "" && val != "" {
				vals[key] = val
			}
		}
	}
	return vals
}

var envs = compositeEnv{belugaEnv("")}

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
