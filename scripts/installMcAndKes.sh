#!/usr/bin/env bash

# install mc (minio client)
curl https://dl.min.io/client/mc/release/linux-amd64/mc \
  --create-dirs \
  -o $HOME/minio-binaries/mc

chmod +x $HOME/minio-binaries/mc

mc --help

# install kes (key management server)
curl -sSL --tlsv1.2 'https://github.com/minio/kes/releases/latest/download/kes-linux-amd64' -o $HOME/minio-binaries/kes
chmod +x $HOME/minio-binaries/kes

export PATH=$PATH:$HOME/minio-binaries/
