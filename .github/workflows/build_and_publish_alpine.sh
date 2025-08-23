#!/bin/sh
set -e

echo "$1" >>/tmp/gnuplex-deploy

apk add uv git go curl bash

export GIT_SSH_COMMAND="ssh -i /tmp/gnuplex-deploy"
git remote set-url origin git@github.com:janie314/gnuplex.git
git config user.name "release workflow"
git config user.email "x@example.com"
git checkout release-linux-musl-x86_64
git merge main -X theirs --allow-unrelated-histories
cd gnuplex
uv run make.py build
git add .
git commit -m "$version_output"
git push origin release-linux-musl-x86_64