#!/bin/bash


MOCK_DIR="./mocks/"

mkdir -p $MOCK_DIR

mockgen -source=./internal/persistence/tradeVariety_repository.go -destination=$MOCK_DIR/tradeVariety/tradeVariety_repository_mock.go -package=mocks
