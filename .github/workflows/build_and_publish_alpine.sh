#!/bin/sh
set -e

umask 077
echo "$1" >/tmp/gnuplex-deploy

apk add git curl bash openssh

curl -LsSf https://astral.sh/uv/install.sh | sh
curl -fsSL https://bun.sh/install | bash
export BUN_INSTALL="$HOME/.bun"
export PATH="$BUN_INSTALL/bin:$HOME/.local/bin:$PATH"

cd /usr/local/bin
wget https://go.dev/dl/go1.25.1.linux-amd64.tar.gz
tar xvzf go*.gz
mv go godir
ln -s go/bin/go .
ln -s go/bin/gofmt .

export GIT_SSH_COMMAND="ssh -i /tmp/gnuplex-deploy -o StrictHostKeyChecking=no"

cd /
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

exit
git push origin release-linux-musl-x86_64