# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM ubuntu:trusty

MAINTAINER Gregor Ambrozic <ambrozic@gmail.com>

RUN apt-get update
RUN apt-get install -y curl

ENV GOVERSION  go1.3.3.linux-amd64
ENV GOROOT /go
ENV PATH $PATH:$GOROOT/bin

RUN curl -O https://storage.googleapis.com/golang/$GOVERSION.tar.gz
RUN tar -C / -xzf $GOVERSION.tar.gz

ADD . /go/src/pkg/mokku
RUN go install mokku

EXPOSE 11222
EXPOSE 8080
