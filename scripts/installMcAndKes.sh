#!/usr/bin/env bash

DIRPATH=$(/usr/bin/dirname $(pwd))
echo $DIRPATH

# install mc (minio client)

curl https://dl.min.io/client/mc/release/linux-amd64/mc \
  --create-dirs \
  -o "$DIRPATH"/minio-binaries/mc
chmod +x "$DIRPATH"/minio-binaries/mc

# install kes (key management server)

curl -sSL --tlsv1.2 'https://github.com/minio/kes/releases/latest/download/kes-linux-amd64' -o "$DIRPATH"/minio-binaries/kes
chmod +x "$DIRPATH"/minio-binaries/kes

# create encryption key for sse-kms encryption

curl 'https://raw.githubusercontent.com/minio/kes/master/root.key' --create-dirs -o "$DIRPATH"/kms-identity/root.key
curl 'https://raw.githubusercontent.com/minio/kes/master/root.cert' --create-dirs -o "$DIRPATH"/kms-identity/root.cert
