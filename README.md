# Programming Challenge: file upload

## Description

The aim of this challenge is to build an application that exposes an API accepting and serving
files. Please document your choices and feel free to ask questions or intermediate code review.

## Preparation

A docker-compose is available to start a local instance of Minio (https://github.com/minio/minio):

```yaml
version: "3"
services:
  minio:
    image: quay.io/minio/minio
    command:
      - server
      - /data
      - --console-address
      - :9001
    ports:
      - "9000:9000"
      - "9001:9001"
```

After starting the service (docker-compose up -d), it will be available from 127.0.0.1:9000 with the
default credentials minioadmin:minioadmin.

## Implementation

The application is written in Go and should:

1. accept files of arbitrary size, encrypt their content, upload them to a Minio bucket
2. serve the submitted files. Submitting and reading files from the API should be possible
   using curl.
3. (bonus) upload the file in chunks of configurable size (e.g. 1MB) in the Minio bucket.
   This mode should be enabled via a configuration. ([multipart file upload](https://github.com/minio/minio-go/blob/master/api-put-object-multipart.go))

## Delivery

Please send us a link to a public Github repository containing your deliverables within 15 days
after receipt of the challenge.

## Solution

### Run

1. Copy `./config/example-config.yml` into `./config/dev-config.yml` and put your own values

#### Run Minio Docker Container

```bash
cd docker/minio
docker-compose up -d
```

Create `access-key-id` and `secret-access-key` and copy the values inside `./config/dev-config.yml`

Create `encryption-key-id` with:

```bash
kes key create dev-key
```

and copy the values inside `./config/dev-config.yml`

**_NOTE:_** to install and configure kes [install mc, kes and configure server side encryption](#kes)

#### Run Go Application

```bash
SERVER_ENVIRONMENT="dev" go run main.go
```

### Run Test Suite

1. With the running minio docker container (verify with `docker ps`)

2. Create small and big files (inside project root folder):

```bash
mkdir testfiles
dd if=/dev/urandom of=./testfiles/small1MiB bs=1M count=1
dd if=/dev/urandom of=./testfiles/big100MiB bs=1M count=100
```

3. Run test suite (inside project root folder)

```bash
go test
```

### cURL calls

#### Upload Big File

```bash
curl --location --request POST 'http://localhost:8080/files' \
--header 'Content-Type: application/json' \
--data-raw '{
    "bucketName":"test",
    "objectName": "big",
    "filepath": "./testfiles/big100MiB",
    "contentType": "application/octet-stream"
}'
```

#### Upload Small File

```bash
curl --location --request POST 'http://localhost:8080/files' \
--header 'Content-Type: application/json' \
--data-raw '{
    "bucketName":"test",
    "objectName": "small",
    "filepath": "./testfiles/small1MiB",
    "contentType": "application/octet-stream"
}'
```

#### Download File (e.g. Small file)

```bash
curl --location --request GET 'http://localhost:8080/files/small' \
--header 'Content-Type: application/json' \
--data-raw '{
    "bucketName": "test",
    "downloadPath": "/tmp/download1.txt"
}'
```

### Multipart Upload

By default minio starts doing multipart upload at 16MiB
The minimum part size is 5MiB
To enable multipart upload there are 2 parameters in `./config/dev-config.yml`:

1. enable-multipart-upload: true
2. file-chunk-size: 5

### Enable Server-Side Encryption (SSE)

<a name="kes"></a>

#### Install mc command (Minio Client)

Source: [https://min.io/docs/minio/linux/reference/minio-mc.html#mc-install](https://min.io/docs/minio/linux/reference/minio-mc.html#mc-install)

```bash
curl https://dl.min.io/client/mc/release/linux-amd64/mc \
  --create-dirs \
  -o $HOME/minio-binaries/mc

chmod +x $HOME/minio-binaries/mc
export PATH=$PATH:$HOME/minio-binaries/

mc --help
```

##### Set alias

```bash
mc alias set myminio http://localhost:9000/ ACCESS_KEY SECRET_KEY
```

Check

```bash
mc admin info myminio
```

#### Install kes command

##### Install kes

Source: [https://github.com/minio/kes](https://github.com/minio/kes)

```bash
curl -sSL --tlsv1.2 'https://github.com/minio/kes/releases/latest/download/kes-linux-amd64' -o $HOME/minio-binaries/kes
chmod +x $HOME/minio-binaries/kes
```

#### 1) Create an Encryption Key for SSE-KMS Encryption

```bash
curl -sSL --tlsv1.2 \
  -O 'https://raw.githubusercontent.com/minio/kes/master/root.key' \
  -O 'https://raw.githubusercontent.com/minio/kes/master/root.cert'
export KES_CLIENT_KEY=kms-identity/root.key
export KES_CLIENT_CERT=kms-identity/root.cert
```

#### 2) Point to play instance

```bash
export KES_SERVER=https://play.min.io:7373
```

##### Create a new EK through KES

```bash
kes key create dev-key
```

This tutorial uses the example my-minio-sse-kms-key name for ease of reference. Specify a unique key name to prevent collision with existing keys. (e.g. dev-key)

#### 2) Configure MinIO for SSE-KMS Object Encryption

Specify the following environment variables in the shell or terminal on each MinIO server host in the deployment (add in the ./docker/minio/docker-compose.yml)

```bash
export MINIO_KMS_KES_ENDPOINT=https://play.min.io:7373
export MINIO_KMS_KES_KEY_FILE=root.key
export MINIO_KMS_KES_CERT_FILE=root.cert
export MINIO_KMS_KES_KEY_NAME=dev-key
```

#### 3) Restart the MinIO Deployment to Enable SSE-KMS

```bash
mc admin service restart myminio
```

#### 4) Configure Automatic Bucket Encryption (Enable Encryption for the bucket)

```bash
mc encrypt set sse-kms dev-key myminio/testbucket
```

**_NOTE:_** to remove automatic encryption on a bucket:

```bash
mc encrypt clear myminio/testbucket
```

**_NOTE:_** And then leave this command hanging in a window to show output logs.

```bash
mc admin trace -v myminio
```

##### Getting Started Running KES Server (Key Encryption Server)

###### 1. Generate KES Server Private Key & Certificate

```bash
kes identity new --ip "127.0.0.1" localhost --cert public.crt --key private.key
```

API Key (secret): kes:v1:AMZWKiawHEsCn30r9t8CV3MvdXVHQ6u4hJvu9q30X0iJ
Identity: 6d2768043cc2350801ef0fbf9d5854541756841470fefe0d517067d1de5f426a

The identity can be computed again via:

```bash
kes identity of kes:v1:AMZWKiawHEsCn30r9t8CV3MvdXVHQ6u4hJvu9q30X0iJ
kes identity of public.crt
```

###### 2. Generate Client Credentials

```bash
kes identity new --key=client.key --cert=client.crt MyApp
```

API Key (secret): kes:v1:AG8LgxNEBAL4aZZvDYzEDA4NKbyog5GMkpvrBgDyUCTr
Identity: dfeef4d0a2b99cf65fb9ebe61275ae341643ef38cb108738d92ba94d36a81f6c

The identity can be computed again via:

```bash
kes identity of kes:v1:AG8LgxNEBAL4aZZvDYzEDA4NKbyog5GMkpvrBgDyUCTr
kes identity of client.crt
```

###### 3. Configure KES Server

Next, we can create the KES server configuration file: config.yml. Please, make sure that the identity in the policy section matches your client.crt identity.

```yaml
address: 0.0.0.0:7373 # Listen on all network interfaces on port 7373

admin:
  identity: dfeef4d0a2b99cf65fb9ebe61275ae341643ef38cb108738d92ba94d36a81f6c # The client.crt identity

tls:
  key: private.key # The KES server TLS private key
  cert: public.crt # The KES server TLS certificate
```

###### 4. Start KES Server

Now, we can start a KES server instance:

```bash
kes server --config config.yml --auth off
```

##### Quickstart KES Server

Source: [https://github.com/minio/kes?tab=readme-ov-file#quick-start](https://github.com/minio/kes?tab=readme-ov-file#quick-start)

1. Configure CLI

```bash
export KES_SERVER=https://play.min.io:7373
export KES_API_KEY=kes:v1:AD9E7FSYWrMD+VjhI6q545cYT9YOyFxZb7UnjEepYDRc
```

2. Create a Key

```bash
kes key create my-key
```

3. Generate a DEK
   Derive a new Data Encryption Keys (DEK)

```bash
kes key dek my-key
```
