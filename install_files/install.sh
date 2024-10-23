#!/bin/sh
printf "Install dir (default $HOME/.local/bin): "
read install_dir
if [ "$install_dir" = "" ]
then
  install_dir="$HOME/.local/bin"
fi
mkdir -p "$install_dir"
git clone -b origin/release https://github.com/janie314/corolla "$install_dir/gnuplex"