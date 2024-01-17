package storage

import (
	"context"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/sse"
	"github.com/pavva91/file-upload/config"
	"github.com/pavva91/file-upload/types"
)

type Storage interface {
	Get(int) *types.User
}

var (
	MinioClient *minio.Client
)

func CreateMinioClient() *minio.Client {
	endpoint := config.ServerConfigValues.Minio.Endpoint
	log.Print(endpoint)

	accessKeyID := config.ServerConfigValues.Minio.AccessKeyID
	secretAccessKey := config.ServerConfigValues.Minio.SecretAccessKey

	useSSL := false
	encryptedBucket := "testbucket"

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	// Set default encryption configuration on a bucket
	err = minioClient.SetBucketEncryption(context.Background(), encryptedBucket, sse.NewConfigurationSSES3())
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("%#v\n", minioClient) // minioClient is now set up
	return minioClient
}
