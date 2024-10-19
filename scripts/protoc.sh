#!/bin/sh

DIR=generated

if [ ! -d "$DIR" ]; then
  mkdir -p $DIR
fi

protoc \
    -I $GOPATH/pkg/mod/github.com/envoyproxy/protoc-gen-validate@v1.0.4 \
    --go_opt=paths=source_relative \
    --go_out=${DIR} \
    --go-grpc_opt=paths=source_relative \
    --go-grpc_out=${DIR} \
    --proto_path=./proto \
    --validate_out=paths=source_relative,lang=go:${DIR} \
    proto/**/*.proto