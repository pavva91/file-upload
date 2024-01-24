#!/usr/bin/env bash

DIRPATH=$(/usr/bin/dirname $(pwd))

export PATH=$PATH:$DIRPATH/minio-binaries

export KES_CLIENT_KEY=$DIRPATH/kms-identity/root.key
export KES_CLIENT_CERT=$DIRPATH/kms-identity/root.cert
export KES_SERVER=https://play.min.io:7373
