#!/usr/bin/env bash

version="0.1"
src="walrustf-0.1"
archive="v$version.tar.gz"

cd /tmp
curl -OL "https://github.com/rcgoodfellow/walrustf/archive/$archive"
tar xzf "$archive"
cd $src

mkdir -p /usr/local/lib/wtf
cp perl/Walrus.pm /usr/local/lib/wtf/
cp python/walrus.py /usr/local/lib/wtf/

#TODO detect C and Go development runtimes and build and install
#     if they are present
