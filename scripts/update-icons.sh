#!/usr/bin/env bash
set -e

echo "Fetching latest nvim-web-devicons..."

VENDOR_DIR="generator/vendor/nvim-web-devicons"
mkdir -p "$VENDOR_DIR"

BASE_URL="https://raw.githubusercontent.com/nvim-tree/nvim-web-devicons/master/lua/nvim-web-devicons/default"

curl -sL "$BASE_URL/icons_by_filename.lua" -o "$VENDOR_DIR/icons_by_filename.lua"
curl -sL "$BASE_URL/icons_by_file_extension.lua" -o "$VENDOR_DIR/icons_by_file_extension.lua"
curl -sL "$BASE_URL/icons_by_operating_system.lua" -o "$VENDOR_DIR/icons_by_operating_system.lua"

echo "Done fetching icons."
