#!/bin/bash

dlv debug cmd/main/main.go --headless --listen=:2345 --api-version=2 --log