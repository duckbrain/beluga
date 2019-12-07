- build

  - ```bash
    docker build -t $BELUGA_IMAGE_TAG $BELUGA_DOCKERFILE $BELUGA_CONTEXT
    docker push -t $BELUGA_IMAGE_TAG
    ```

- env

  - Produces variables in `.env`, `bash`, etc formats to be used in scripts

- deploy

  - prep compose file

    ```bash
    COMPOSE_FILE=$BELUGA_COMPOSE_FILE docker-compose config > /gen/docker-compose.yaml
    ```

  - transport: (swarm/compose) type determined by `DOCKER_HOST`

    - local (default)
    - ssh `DOCKER_HOST=ssh://username@hostname`
      - SSH forward the socket, run swarm/compose in the image (need standard directory path)
    - portainer (perfered)
    - Only one that doesn't run `docker` or `docker-compose` directly.
  
- portainer (option 0)
  
    ```bash
    # pretend it's authenticated
    curl 
  ```
  
- swarm (option 1)
  
    ```bash
    docker stack deploy \
    	--compose-file /gen/docker-compose.yaml \
    	--prune \
    	--with-registry-auth
  ```
  
- docker-compose (option 2)
  
    ```bash
    docker-compose up \
    	--file /gen/docker-compose.yaml \
    	--no-build \
    	--remove-orphans \
    	--detach
    ```
  
- teardown: same as deploy

- 