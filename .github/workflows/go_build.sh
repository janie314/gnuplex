#!/bin/sh
set -e
apk add uv git go curl bash
git clone --depth 1 --branch 2025.08.10.jd-alpine-install-script https://github.com/janie314/gnuplex
cd gnuplex
uv run make.py go_build