# Base image: https://hub.docker.com/_/golang/
FROM golang:1.16.3
USER root

ENV GOPATH /go
ENV PATH ${GOPATH}/bin:$PATH

ENV GOPROXY=https://goproxy.io,direct
ENV GO111MODULE=on

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' > /etc/timezone

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.39.0
RUN curl -sfL https://raw.githubusercontent.com/reviewdog/reviewdog/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v0.11.0