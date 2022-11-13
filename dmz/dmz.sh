#!/usr/bin/sh
basedir=$(realpath $(dirname "$0"))
tinyproxy -d -c "$basedir/dmz.conf"
