# Build mola docker image

## Get source code
1. git clone https://github.com/Dataman-Cloud/log-proxy.git

## Build
1. build mola binaray file on CentOS 7 and Go 1.7.

```
# cd log-proxy
# docker run -it --rm --name buildmola -v "$PWD":/go/src/github.com/Dataman-Cloud/log-proxy -w /go/src/github.com/Dataman-Cloud/log-proxy golang:1.7.4-onbuild make build
# ls bin/mola
```

2. compress the UI files

```
# cd frotend
# docker run -it --rm --name buildmola -v "$PWD":/usr/src/app -w /usr/src/app node:4.7.0-onbuild ./compress.sh
```

3. build the image  

```
docker build -f dockerfiles/Dockerfile_runtime -t mola:0.1 .
```

## Run

```
docker run -d -p 5098:5098 -e ES_URL=http://127.0.0.1:9200 -e PROMETHEUS_URL=http://127.0.0.1:9090 -e ALERTMANAGER_URL=http://127.0.0.1:9093 --name mola mola:0.1
```
