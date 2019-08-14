package lib

const (
	varApplication      = "BELUGA_APPLICATION"
	varDockerContext    = "BELUGA_DOCKER_CONTEXT"
	varDockerfile       = "BELUGA_DOCKERFILE"
	varDomain           = "BELUGA_DOMAIN"
	varEnvironment      = "BELUGA_ENVIRONMENT"
	varGitDefaultBranch = "BELUGA_GIT_DEFAULT_BRANCH"
	varImage            = "BELUGA_IMAGE"
	varRegistry         = "BELUGA_REGISTRY"
	varRegistryPassword = "BELUGA_REGISTRY_PASSWORD"
	varRegistryUsername = "BELUGA_REGISTRY_USERNAME"
	varVariant          = "BELUGA_VARIANT"
	varVersion          = "BELUGA_VERSION"
)

var knownVarNames = []string{
	varApplication,
	varDockerContext,
	varDockerfile,
	varDomain,
	varEnvironment,
	varGitDefaultBranch,
	varImage,
	varRegistry,
	varRegistryPassword,
	varRegistryUsername,
	varVariant,
	varVersion,
}

func (e Env) Application() string {
	return e[varApplication]
}

func (e Env) DockerContext() string {
	return e[varDockerContext]
}

func (e Env) Dockerfile() string {
	return e[varDockerfile]
}

func (e Env) Domain() string {
	return e[varDomain]
}

func (e Env) Environment() string {
	return e[varEnvironment]
}

func (e Env) GitDefaultBranch() string {
	return e[varGitDefaultBranch]
}

func (e Env) Image() string {
	return e[varImage]
}

func (e Env) Registry() string {
	return e[varRegistry]
}

func (e Env) RegistryPassword() string {
	return e[varRegistryPassword]
}

func (e Env) RegistryUsername() string {
	return e[varRegistryUsername]
}

func (e Env) Variant() string {
	return e[varVariant]
}

func (e Env) Version() string {
	return e[varVersion]
}
