#!/bin/bash


MOCK_DIR="./mocks/"

mkdir -p $MOCK_DIR

mockgen -source=./internal/persistence/tradeVariety_repository.go -destination=$MOCK_DIR/persistence/trade_variety/tradeVariety_repository_mock.go -package=mocks
mockgen -source=./internal/persistence/asset_repository.go -destination=$MOCK_DIR/persistence/asset/asset_repository_mock.go -package=mocks