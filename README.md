#mokku

Monitoring service for docker cluster.

Server collects metrics from client instances running on servers in docker cluster. Data is collected using [control groups](https://www.kernel.org/doc/Documentation/cgroups/cgroups.txt). For now only collected metrics are about OS, CPU and memory.

##Docker
```shell
docker pull ambrozic/mokku
```

##Server
```shell
docker rm -f mokku-server || true &&
docker run \
    --name mokku-server \
    -i \
    -t \
    --publish 11222:11222 \
    --publish 8080:8080 \
    ambrozic/mokku \
    /go/bin/mokku --server
```
####Example
```shell
docker rm -f mokku-server || true && docker run --name mokku-server -i -t -p 11222:11222 -p 8080:8080 ambrozic/mokku /go/bin/mokku --server
```

##Client
```shell
docker rm -f mokku || true &&
docker run \
    --name mokku \
    -i \
    -t \
    --env DOCKER_HOST=tcp://server-ip:port \
    --volume /sys/fs/cgroup:/sys/fs/cgroup:ro \
    ambrozic/mokku \
    /go/bin/mokku --host=server-ip --port=11222
```
####Example
```shell
docker rm -f mokku || true && docker run --name mokku -i -t -e DOCKER_HOST=tcp://192.168.59.103:2375 -v /sys/fs/cgroup:/sys/fs/cgroup:ro ambrozic/mokku /go/bin/mokku --host=192.168.59.103 --port=11222
```
