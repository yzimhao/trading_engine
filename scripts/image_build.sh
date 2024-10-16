#!/bin/bash

echo $(pwd)

TAG=$1
if [ "X${TAG}" = "X" ];then
    echo "image tag cannot be empty, example: ./build-image.sh v1.5.8"
    exit 1
fi

echo "docker build -f Dockerfile -t yzimhao/haotrader:${TAG} .."
docker build -f Dockerfile -t yzimhao/haotrader:${TAG} ..