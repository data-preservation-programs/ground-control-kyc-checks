# syntax=docker/dockerfile:1

# https://docs.docker.com/language/golang/build-images/

FROM golang:1.18-alpine

RUN apk add build-base
RUN apk add jq

WORKDIR /usr/src/app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY testdata ./testdata
COPY checks ./checks

RUN go build -o sp-kyc-checks cmd/main.go 

RUN GOOGLE_MAPS_API_KEY=skip MAXMIND_USER_ID=skip go test ./checks/geoip

CMD /usr/src/app/sp-kyc-checks testdata/responses-1-pass.json | tee test-results.json
