#!/bin/sh
set -eu

export PATH="$PATH:/usr/local/go/bin"
export CC="gcc -std=gnu99"

go build -o gslauncher-linux-i386 ../cmd/gslauncher
go build -tags debug -o gslauncher-linux-i386-debug ../cmd/gslauncher
