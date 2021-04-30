#!/bin/sh
set -eu

if [ "$(id -u)" -ne 0 ]; then
    echo 'root privileges required' >&2
    exit 1
fi

binname=

case "$(uname -m)" in
    x86_64)
        binname=gslauncher-linux-amd64
        ;;
    i386|i686)
        binname=gslauncher-linux-i386
        ;;
    *)
        echo 'platform not supported' >&2
        exit 1
        ;;
esac

mkdir /opt/gslauncher
install -m 755 -o root -g root $binname /opt/gslauncher
install -m 644 -o root -g root gslauncher.desktop /opt/gslauncher
install -m 644 -o root -g root logo.png /opt/gslauncher
ln -s /opt/gslauncher/gslauncher.desktop /usr/share/applications/
