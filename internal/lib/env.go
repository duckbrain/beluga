package lib

import (
	"github.com/gobuffalo/envy"
)

const (
	envDSN               = "DSN"
	envDomain            = "DOMAIN"
	envDockerUsername    = "DOCKER_USERNAME"
	envDockerPassword    = "DOCKER_PASSWORD"
	envDockerImage       = "DOCKER_IMAGE"
	envGitRefName        = "GIT_REF_NAME"
	envGitDefaultRefName = "GIT_DEFAULT_REF_NAME"
)

type Env map[string]string

func (e Env) DSN() string {
	return e[envDSN]
}

func (e Env) Domain() string {
	return e[envDomain]
}

func (e Env) DockerCreds() (username, password string) {
	username = e[envDockerUsername]
	password = e[envDockerPassword]
	return
}

func (e Env) Client() Client {
	return Client{
		DSN: e[envDSN],
	}
}

func DeploymentFromEnv(e Env) (Deployment, err) {

}

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

var envs = compositeEnv{belugaEnv("BELUGA_")}

type belugaEnv string

func (prefix belugaEnv) EnvRead() Env {
	vals := make(Env)
	keys := []string{
		envDockerUsername, envDockerPassword, envDockerImage,
		envGitRefName, envGitDefaultRefName,
		envDomain,
	}
	for _, key := range keys {
		val := envy.Get(string(prefix)+key, "")
		if len(val) > 0 {
			vals[key] = val
		}
	}
	return vals
}
