#!/bin/bash
set -e

OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

BIN="upstat-linux"
if [ "$OS" == "darwin" ]; then
    BIN="upstat-mac"
elif [ "$OS" == "linux" ]; then
    BIN="upstat-linux"
else
    echo "Sistema não suportado"
    exit 1
fi

curl -L -o /tmp/$BIN https://github.com/yagofontanez/upstat-cli-go/releases/latest/download/$BIN
chmod +x /tmp/$BIN
sudo mv /tmp/$BIN /usr/local/bin/upstat

echo "UpStat CLI instalado! Rode 'upstat' para usar."