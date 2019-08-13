package lib

const (
	varApplication      = "BELUGA_APPLICATION"
	varDomain           = "BELUGA_DOMAIN"
	varEnvironment      = "BELUGA_ENVIRONMENT"
	varImage            = "BELUGA_IMAGE"
	varRegistry         = "BELUGA_REGISTRY"
	varRegistryPassword = "BELUGA_REGISTRY_PASSWORD"
	varRegistryUsername = "BELUGA_REGISTRY_USERNAME"
	varVariant          = "BELUGA_VARIANT"
	varVersion          = "BELUGA_VERSION"
)

var knownVarNames = []string{
	varApplication,
	varDomain,
	varEnvironment,
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

func (e Env) Domain() string {
	return e[varDomain]
}

func (e Env) Environment() string {
	return e[varEnvironment]
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
