#!/bin/sh
set -eu

cd "$(dirname $0)/.."

export PATH="$PATH:/usr/local/go/bin"
export CGO_ENABLED=1
export GOOS=windows

export CC=x86_64-w64-mingw32-gcc
export GOARCH=amd64

x86_64-w64-mingw32-windres -o ./cmd/gslauncher/logo.syso build/logo.rc
go build -ldflags '-H=windowsgui' -o dist/gslauncher-windows-amd64.exe ./cmd/gslauncher
go build -tags debug -o dist/gslauncher-windows-amd64-debug.exe ./cmd/gslauncher


export CC=i686-w64-mingw32-gcc
export GOARCH=386

i686-w64-mingw32-windres -o ./cmd/gslauncher/logo.syso build/logo.rc
go build -ldflags '-H=windowsgui' -o dist/gslauncher-windows-i386.exe ./cmd/gslauncher
go build -tags debug -o dist/gslauncher-windows-i386-debug.exe ./cmd/gslauncher


rm ./cmd/gslauncher/logo.syso


makensis build/installer.nsi
