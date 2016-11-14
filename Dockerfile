FROM alpine:3.4

ENV GO15VENDOREXPERIMENT 1
ENV GOPATH /go

COPY . /go/src/github.com/Dataman-Cloud/log-proxy/
RUN cd /go/src/github.com/Dataman-Cloud/log-proxy/ && ./build.sh 

WORKDIR /bin

ENTRYPOINT ./log-proxy
CMD ["-config", "env_file"]
