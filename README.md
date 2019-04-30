# Beluga

A tool for building and deploying docker web applications in CI without Kubernetes.


## `beluga build`

Pulls the latest build, then uses it as a cache, to build and push a docker image to a registry.

## `beluga beach [COMPOSE_FILE]`

Sends the `docker-compose.yaml` file to the belugad server, where it will be loaded.

## `belugad serve`

Starts the deamon to run with nginx-proxy and, optionally, nginx-proxy-companion. The only time you'll need to interact with this is to set it up.


