#!/bin/sh
set -e

umask 077
echo "$1" >/tmp/gnuplex-deploy

apk add uv git go curl bash openssh

curl -fsSL https://bun.sh/install | bash
export BUN_INSTALL="$HOME/.bun"
export PATH="$BUN_INSTALL/bin:$PATH"

export GIT_SSH_COMMAND="ssh -i /tmp/gnuplex-deploy -o StrictHostKeyChecking=no"

git clone git@github.com:janie314/gnuplex.git
cd gnuplex
git config user.name "release workflow"
git config user.email "x@example.com"
git checkout release-linux-musl-x86_64
git pull origin release-linux-musl-x86_64

git merge main -X theirs --allow-unrelated-histories

uv run make.py build
version_output=$(./backend/bin/gnuplex -version)
git add .
git commit -m "$version_output"

git push origin release-linux-musl-x86_64