#!/bin/sh
printf "Install dir (default $HOME/.local/bin): "
read install_dir
if [ "$install_dir" = "" ]
then
  install_dir="$HOME/.local/bin"
fi
mkdir -p "$install_dir"
git clone -b release https://github.com/janie314/gnuplex.git "$install_dir/gnuplex-code"
