#!/bin/sh
set -eux

cd "$(dirname $0)"

docker build -f Dockerfile-linux-amd64 . -t gslauncher-linux-build:amd64
docker buildx build --platform i386 -f Dockerfile-linux-i386 . -t gslauncher-linux-build:i386
docker build -f Dockerfile-cross-win . -t gslauncher-win-build

docker run -v $(pwd)/..:/data gslauncher-linux-build:amd64 /data/build/build-linux.sh amd64
docker run --platform i386 -v $(pwd)/..:/data gslauncher-linux-build:i386 /data/build/build-linux.sh i386

TMPDIR="$(mktemp -d)"
mkdir "${TMPDIR}/gslauncher"
install -m 755 ../dist/gslauncher-linux-{amd64,i386} setup.sh "${TMPDIR}/gslauncher/"
install -m 644 gslauncher.desktop logo.png "${TMPDIR}/gslauncher/"
tar -C "${TMPDIR}" -czf ../dist/gslauncher-linux.tar.gz gslauncher
rm -r $TMPDIR

docker run -v $(pwd)/..:/data gslauncher-win-build /data/build/build-win.sh
