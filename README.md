# DIST-APP
Sample implementation of gossip protocol

## lint
```bash
$ brew install golangci-lint 
$ make lint
```
or
```bash
$ make go-lint-docker
```

## Build
### Mac
```bash
$ make mac_amd64_build
```

### Linux
```bash
$ make linux_build
```

### Docker Build
```bash
$ make linux_build

$ make docker_build
# using Dockerfile docker image will be build 
# and saved locally as distapp:0.0.1
```

## Docker Deployment
```bash
$ make deploy 
# Note: this will deploy 3 instance (prefix with distapp*)

$ make destroy 
# Note: will stop and remove deployed 3 instances 
# (prefix with distapp*)
```

## Cockroach DB deployment
Step 1. Create a bridge network
```bash
$ docker network create -d bridge roachnet
```

Step 2. Start the cluster
1. Create a Docker volume for each container:
```bash
$ docker volume create roach1
$ docker volume create roach2
$ docker volume create roach3
```

2. Start the first node:
```bash
$ docker run -d \
--name=roach1 \
--hostname=roach1 \
--net=roachnet \
-p 26257:26257 -p 8080:8080  \
-v "roach1:/cockroach/cockroach-data"  \
cockroachdb/cockroach:v22.1.6 start \
--insecure \
--join=roach1,roach2,roach3
```

3. Start two more nodes:
```bash
$ docker run -d \
--name=roach2 \
--hostname=roach2 \
--net=roachnet \
-v "roach2:/cockroach/cockroach-data" \
cockroachdb/cockroach:v22.1.6 start \
--insecure \
--join=roach1,roach2,roach3

$ docker run -d \
--name=roach3 \
--hostname=roach3 \
--net=roachnet \
-v "roach3:/cockroach/cockroach-data" \
cockroachdb/cockroach:v22.1.6 start \
--insecure \
--join=roach1,roach2,roach3
```

```bash
$ docker exec -it roach1 ./cockroach init --insecure
```