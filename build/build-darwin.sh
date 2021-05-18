#!/bin/sh
set -eux

cd "$(dirname $0)/.."
SRCDIR="$(pwd)"

go install fyne.io/fyne/v2/cmd/fyne@latest

TMPDIR="$(mktemp -d)"
cd $TMPDIR
fyne package -release -icon "${SRCDIR}/build/logo.png" -sourceDir "${SRCDIR}/cmd/gslauncher" -name 'GrooveStats Launcher'
ln -s /Applications Applications
cd -
hdiutil create dist/gslauncher-macos-amd64.dmg -fs HFS+ -volname 'GrooveStats Launcher' -format UDZO -srcfolder $TMPDIR
rm -r $TMPDIR

TMPDIR="$(mktemp -d)"
cd $TMPDIR
fyne package -tags debug -icon "${SRCDIR}/build/logo.png" -sourceDir "${SRCDIR}/cmd/gslauncher" -name 'GrooveStats Launcher (debug)'
ln -s /Applications Applications
cd -
hdiutil create dist/gslauncher-macos-amd64-debug.dmg -fs HFS+ -volname 'GrooveStats Launcher (debug)' -format UDZO -srcfolder $TMPDIR
rm -r $TMPDIR
