#!/bin/sh
set -e
printf "Install dir (default %s/.local/bin): " "$HOME"
read -r install_dir
if [ "$install_dir" = "" ]
then
  install_dir="$HOME/.local/bin"
fi
mkdir -p "$install_dir"
git clone -b release https://github.com/janie314/gnuplex.git "$install_dir/gnuplex-code"
cd "$install_dir"
ln -s gnuplex-code/backend/bin/gnuplex .