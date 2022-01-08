#!/bin/sh
set -eux

cd "$(dirname $0)/.."
SRCDIR="$(pwd)"

export CGO_ENABLED=1
export GOOS=darwin

export GOARCH=amd64
go build -o dist/gslauncher-darwin-amd64 ./cmd/gslauncher
go build -tags debug -o dist/gslauncher-darwin-amd64-debug ./cmd/gslauncher

export GOARCH=arm64
go build -o dist/gslauncher-darwin-arm64 ./cmd/gslauncher
go build -tags debug -o dist/gslauncher-darwin-arm64-debug ./cmd/gslauncher

unset GOARCH

lipo -create -output dist/gslauncher-darwin dist/gslauncher-darwin-{amd64,arm64}
lipo -create -output dist/gslauncher-darwin-debug dist/gslauncher-darwin-{amd64,arm64}-debug

go install fyne.io/fyne/v2/cmd/fyne@latest

TMPDIR="$(mktemp -d)"
cd $TMPDIR
fyne package -release -executable "${SRCDIR}/dist/gslauncher-darwin" -icon "${SRCDIR}/build/logo.png" -name 'GrooveStats Launcher'
ln -s /Applications Applications
cd -
hdiutil create dist/gslauncher-macos.dmg -fs HFS+ -volname 'GrooveStats Launcher' -format UDZO -srcfolder $TMPDIR
rm -r $TMPDIR

TMPDIR="$(mktemp -d)"
cd $TMPDIR
fyne package -executable "${SRCDIR}/dist/gslauncher-darwin-debug" -icon "${SRCDIR}/build/logo.png" -name 'GrooveStats Launcher (debug)'
ln -s /Applications Applications
cd -
hdiutil create dist/gslauncher-macos-debug.dmg -fs HFS+ -volname 'GrooveStats Launcher (debug)' -format UDZO -srcfolder $TMPDIR
rm -r $TMPDIR
