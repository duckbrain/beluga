# Beluga

A tool for building and deploying web applications in CI through Docker Compose or Swarm.

- GitLab ![Pipeline Status](https://gitlab.com/duckbrain/beluga/badges/master/pipeline.svg)

## Commands

### `beluga build`

Pulls the latest build, then uses it as a cache, to build and push a docker image to a registry.

### `beluga deploy`

Sends the `docker-compose.yaml` file to the belugad server, where it will be loaded.

### `beluga teardown`

Instructs belugad to teardown the stack deployed for a `beluga build` command.

## Environment Variables

## Terminology

- *Stack*: refers to an application stack to be delpoyed on belugad. It is referenced by the domain name.
- *DSN*: is a string for connecting to belugad. It contains the domain name of the belugad instance, and a key for deploying.