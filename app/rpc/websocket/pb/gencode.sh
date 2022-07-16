#!/bin/bash
goctl rpc protoc websocket.proto -v --go_out=../ --go-grpc_out=../  --zrpc_out=../ --style=goZero