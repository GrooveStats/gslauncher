#!/bin/sh
set -eu

export PATH="$PATH:/usr/local/go/bin"
export CGO_ENABLED=1
export GOOS=windows

export CC=x86_64-w64-mingw32-gcc
export GOARCH=amd64

go build -o gslauncher-windows-amd64.exe ../cmd/gslauncher
go build -tags debug -o gslauncher-windows-amd64-debug.exe ../cmd/gslauncher


export CC=i686-w64-mingw32-gcc
export GOARCH=386

go build -o gslauncher-windows-i386.exe ../cmd/gslauncher
go build -tags debug -o gslauncher-windows-i386-debug.exe ../cmd/gslauncher
