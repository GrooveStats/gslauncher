#!/bin/sh
set -eu

if [ $# -ne 1 ]; then
	echo "usage: $0 amd64|i386" >&2
	exit 1
fi

arch="$1"

cd "$(dirname $0)/.."

export PATH="$PATH:/usr/local/go/bin"
export CC="gcc -std=gnu99"

go build -o "dist/gslauncher-linux-${arch}" ./cmd/gslauncher
go build -tags debug -o "dist/gslauncher-linux-${arch}-debug" ./cmd/gslauncher
