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
   This mode should be enabled via a configuration.

## Delivery

Please send us a link to a public Github repository containing your deliverables within 15 days
after receipt of the challenge.
