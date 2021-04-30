#!/bin/sh
set -eux

docker build -f Dockerfile-linux-amd64 . -t gslauncher-linux-build:amd64
docker buildx build --platform i386 -f Dockerfile-linux-i386 . -t gslauncher-linux-build:i386
docker build -f Dockerfile-cross-win . -t gslauncher-win-build

docker run -v $(pwd)/..:/data -w /data/build gslauncher-linux-build:amd64 ./build-linux-amd64.sh
docker run --platform i386 -v $(pwd)/..:/data -w /data/build gslauncher-linux-build:i386 ./build-linux-i386.sh
docker run -v $(pwd)/..:/data -w /data/build gslauncher-win-build ./build-win.sh

TMPDIR="$(mktemp -d)"
mkdir "${TMPDIR}/gslauncher"
install -m 755 gslauncher-linux-i386 gslauncher-linux-amd64 setup.sh "${TMPDIR}/gslauncher/"
install -m 644 gslauncher.desktop logo.png "${TMPDIR}/gslauncher/"
tar -C "${TMPDIR}" -czf gslauncher-linux.tar.gz gslauncher
rm -r $TMPDIR
