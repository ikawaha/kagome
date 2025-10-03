# Example of using Kagome's Docker Image with Graphviz support

If you are familiar with `docker` and `docker compose`, you can easily set up a Kagome server with [Graphviz](https://graphviz.org/) support for visualizing lattice structures.

> [!NOTE]
> The [docker image we provide](https://hub.docker.com/r/ikawaha/kagome) does not include Graphviz due to size constraints. This example shows how to build a custom docker image that includes Graphviz.

* [Dockerfile](./Dockerfile)
* [docker-compose.yml](./docker-compose.yml)
* [userdict.txt](./userdict.txt)

## Bootstrapping the Environment

```shell
docker compose up -d
```

After running the above command, you can access the Kagome server at `http://localhost:6060/`.

## Shutting Down the Environment

```shell
docker compose down
```
