# Beluga

A tool for building and deploying docker web applications in CI without Kubernetes.

## Example Usage

In ideal an ideal situation you can set up the server with HTTP or HTTPS by using one of the `docker-compose.yaml` files below.

```yaml
version: "3.5"
services: # TODO
```

## Commands

### `beluga build`

Pulls the latest build, then uses it as a cache, to build and push a docker image to a registry.

### `beluga deploy [COMPOSE_FILE]`

Sends the `docker-compose.yaml` file to the belugad server, where it will be loaded.

### `beluga teardown`

Instructs belugad to teardown the stack deployed for a `beluga build` command.

## Environment Variables


| Data point                 | GitLab                            | Default                 |
|----------------------------|-----------------------------------|-------------------------|
| Domain Name                | Parsed from `$CI_ENVIRONMENT_URL` |                         |
| `docker-compose.yaml` path | `$BELUGA_COMPOSE`                 | `./docker-compose.yaml` |
| DSN                        | `$BELUGA_DSN`                     |                         |

## Terminology

- *Domain*: is a domain name that should be pointed to the `nginx-proxy` host.
- *Stack*: refers to an application stack to be delpoyed on belugad. It is referenced by the domain name.
- *DSN*: is a string for connecting to belugad. It contains the domain name of the belugad instance, and a key for deploying.