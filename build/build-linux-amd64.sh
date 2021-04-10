#!/bin/sh
set -eu

export PATH="$PATH:/usr/local/go/bin"
export CC="gcc -std=gnu99"

go build -o gslauncher-linux-amd64 ../cmd/gslauncher
go build -tags debug -o gslauncher-linux-amd64-debug ../cmd/gslauncher
