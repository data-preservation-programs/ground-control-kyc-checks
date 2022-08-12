#! /bin/bash

. .env

export GOOGLE_MAPS_API_KEY

go test ./...

