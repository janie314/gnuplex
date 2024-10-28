#!/bin/sh
set -e
# dependencies
echo "checking if git exists"
command -v git
echo "checking if curl exists"
command -v curl
# continue
printf "Install dir (default %s/.local/bin): " "$HOME"
read -r install_dir
if [ "$install_dir" = "" ]
then
  install_dir="$HOME/.local/bin"
fi
mkdir -p "$install_dir"
git clone -b release https://github.com/janie314/gnuplex.git "$install_dir/gnuplex-code"
mkdir -p "$install_dir/gnuplex-code/backend/bin"
curl -o "$install_dir/gnuplex-code/backend/bin/gnuplex" https://gnuplex.janie.page/gnuplex
cd "$install_dir"
ln -s gnuplex-code/backend/bin/gnuplex .
