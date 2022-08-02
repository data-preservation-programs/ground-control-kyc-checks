# syntax=docker/dockerfile:1

# https://docs.docker.com/language/golang/build-images/

FROM golang:1.18-alpine

RUN apk add build-base
RUN apk add jq

WORKDIR /usr/src/app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

COPY testdata ./testdata
COPY minpower ./minpower

RUN go build

WORKDIR /usr/src/app/minpower
RUN go test
WORKDIR /usr/src/app

CMD /usr/src/app/sp-kyc-checks testdata/responses-1-pass.json | tee test-results.json
