#!/bin/bash
# Prints out the directory tree and files of your project

# Stop on errors
set -e

# Starting point
ROOT_DIR=${1:-.}

echo "📂 Directory structure for: $(realpath "$ROOT_DIR")"
echo "----------------------------------------"

# Use find to list files and directories with indentation
find "$ROOT_DIR" -print | sed -e 's;[^/]*/;|____;g;s;____|; |;g'

echo "----------------------------------------"
echo "✅ Done!"
