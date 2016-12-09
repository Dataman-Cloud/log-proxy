#!/bin/sh
apk add --update go gcc g++
go build -v -o ./bin/log-proxy ./src/ && \
	mv bin/log-proxy /bin/log-proxy && \
	mv env_file.template /bin/env_file

apk del go gcc g++
rm -rf /go && rm -rf /var/cache/apk/*
