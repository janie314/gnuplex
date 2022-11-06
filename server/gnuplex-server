#!/usr/bin/env bash

set -o errexit
set -o pipefail
set -o nounset
# set -o xtrace

SOCKET="/tmp/mpvsocket"

if test -f "${SOCKET}"; then
  rm "${SOCKET}"
fi
mpv --idle=yes --input-ipc-server="${SOCKET}" --fs
