#!/bin/bash
goctl api go --api api.api -dir ../  --style=goZero
# shellcheck disable=SC2035
rm -rf *.md
goctl api doc --dir . --o .
