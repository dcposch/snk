#!/bin/sh
set -e
export OS=`uname`
if test "$OS" = "Linux"; then export OS="linux";
elif test "$OS" = "Darwin"; then export OS="mac";
fi
echo "Supported platforms: mac, linux. Found: $OS"
FILE=`mktemp`
curl -sfL https://github.com/dcposch/snk/releases/download/v1.0.0/snk-$OS > $FILE
chmod +x $FILE
$FILE