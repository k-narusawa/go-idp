#!/bin/bash

rm -rf go
rm -rf typescript

docker run --rm \
    -v $PWD/openapi:/local openapitools/openapi-generator-cli:v7.2.0 generate \
    -i /local/src/openapi.yml \
    -g go-echo-server \
    -o /local/go

docker run --rm \
    -v $PWD/openapi:/local openapitools/openapi-generator-cli:v7.2.0 generate \
    -i /local/src/openapi.yml \
    -g typescript-axios \
    --additional-properties platform=node \
    -o /local/typescript