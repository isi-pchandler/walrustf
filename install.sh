#!/usr/bin/env bash

version="0.1"
src="walrus-$version"
archive="$src.tar.gz"

cd /tmp
curl -OL "https://github.com/rcgoodfellow/walrustf/archive/$archive"
tar xzf "$archive"
cd $src

mkdir /usr/local/lib/wtf
cp perl/Walrus.pm /usr/local/lib/wtf/
cp python/Walrus.py /usr/local/lib/wtf/

#TODO detect C and Go development runtimes and build and install
#     if they are present
