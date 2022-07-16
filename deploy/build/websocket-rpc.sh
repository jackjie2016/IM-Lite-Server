#!/bin/bash
# 以下是M1芯片的编译脚本，请根据自己的需求修改
#brew tap messense/macos-cross-toolchains
#brew install x86_64-unknown-linux-gnu
#brew install aarch64-unknown-linux-gnu
cd ../../app/rpc/websocket
CGO_ENABLED=1 GOOS=linux go build --ldflags "-extldflags -static" -o main . || exit 1
cd -
cd websocket-rpc
docker build --platform=linux/amd64 -t registry.cn-shanghai.aliyuncs.com/pathim/websocket-rpc:latest . || exit 1
docker push registry.cn-shanghai.aliyuncs.com/pathim/websocket-rpc:latest || exit 1
