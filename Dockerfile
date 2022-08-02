# syntax=docker/dockerfile:1

# https://docs.docker.com/language/golang/build-images/

FROM golang:1.18-alpine

RUN apk add jq

WORKDIR /usr/src/app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN go build

CMD /usr/local/go/bin/go testdata/responses-1-pass.json | tee test-results.json
