#!/bin/sh
set -e
echo "Installing GNUPlex..."
# dependencies
echo "checking if git exists"
command -v git
echo "checking if curl exists"
command -v curl
# continue
echo "All needed commands exist, continuing..."
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
echo "Installing systemd user service..."
mkdir -p "$HOME/.config/systemd/user"
install_dir_replace=$(echo "$install_dir/gnuplex-code/backend" | sed -e 's~/~\\/~g')
cat "$install_dir/gnuplex-code/install_files/gnuplex.service" | sed -e "s/__DIR__/$install_dir_replace/" >"$HOME/.config/systemd/user/gnuplex.service"
echo "Done. To turn on GNUPlex, run: systemctl --user enable --now gnuplex"