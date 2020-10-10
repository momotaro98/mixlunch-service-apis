# mixlunch-service-api

## Use as a service

### Up docker compose

```bash
$ make docker-run
```

You can stop the containers

```bash
$ make docker-stop
```

## Development

### Set up local environment variables

```
$ direnv allow
```

### Help command in Makefile

You can check `make` commands in `Makefile`

```
$ make help
```

### Setup

First you need to install go tools

```
$ make setup
```

### Build and Run tests

```bash
$ make
```

### Run only middleware containers to run app as local program not Docker container

#### The command below runs containers of DB and cache and so on.

```bash
$ make docker-mid-run
$ make docker-stop
```

#### You can do **Hot Reload** development

```bash
$ make hot-reload
```

### How to update gRPC/protocol buffer

```
$ cd mixlunch-service-api
$ protoc -I pb/ pb/mixlunch.proto --go_out=plugins=grpc:pb
```
