#!/bin/sh
set -e
echo "$1" >>/tmp/gnuplex-deploy
apk add uv git go curl bash
git clone --depth 1 --branch release https://github.com/janie314/gnuplex
cd gnuplex
uv run make.py build
git add .
git commit -m "$version_output"
GIT_SSH_COMMAND="ssh -i /tmp/gnuplex-deploy" git push origin release-linux-musl-x86_64