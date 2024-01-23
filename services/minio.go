// FileUploader.go MinIO example
package services

import (
	"context"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/encrypt"
	"github.com/pavva91/file-upload/config"
	"github.com/pavva91/file-upload/storage"
)

func UploadFile(objectName string, filePath string, contentType string, bucketName string) (minio.UploadInfo, error) {
	ctx := context.Background()

	// Upload the test file
	// Change the value of filePath if the file is in another location

	// Upload the test file with FPutObject
	info, err := storage.MinioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
		return minio.UploadInfo{}, err
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
	return info, nil
}

// explicitly encrypt
// NOTE: I can configure automatic bucket encryption with mc
// mc encrypt set sse-kms dev-key myminio/testbucket
func EncryptAndUploadFile(objectName string, filePath string, contentType string, bucketName string) (minio.UploadInfo, error) {
	ctx := context.Background()

	// encryption, err := encrypt.NewSSEKMS("dev-key2", ctx)
	encryption, err := encrypt.NewSSEKMS(config.ServerConfigValues.Minio.EncryptionKeyID, ctx)
	if err != nil {
		log.Println(err)
		return minio.UploadInfo{}, err
	}

	// Upload the test file with FPutObject
	info, err := storage.MinioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ServerSideEncryption: encryption})
	if err != nil {
		log.Println(err)
		return minio.UploadInfo{}, err
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
	return info, nil
}

// TODO: Must enable SSL
func EncryptWithPasswordAndUploadFile(objectName string, filePath string, contentType string, bucketName string) (minio.UploadInfo, error) {
	ctx := context.Background()

	password := "correct horse battery staple"

	// New SSE-C where the cryptographic key is derived from a password and the objectname + bucketname as salt
	encryption := encrypt.DefaultPBKDF([]byte(password), []byte(bucketName+objectName))

	// Upload the test file with FPutObject
	info, err := storage.MinioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ServerSideEncryption: encryption})
	if err != nil {
		log.Println(err)
		return minio.UploadInfo{}, err
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
	return info, nil
}

func DownloadFile(bucket string, fileName string, downloadPath string) error {
	err := storage.MinioClient.FGetObject(context.Background(), bucket, fileName, downloadPath, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	log.Printf("File %s correctly downloaded in: %s", fileName, downloadPath)
	return nil
}

func BucketExist(bucket string) (bool, error) {
	found, err := storage.MinioClient.BucketExists(context.Background(), bucket)
	if err != nil {
		return false, err
	}
	return found, nil
}

func EncryptAndUploadFileMultipart(objectName string, filePath string, contentType string, bucketName string) (minio.UploadInfo, error) {
	ctx := context.Background()

	// encryption, err := encrypt.NewSSEKMS("dev-key2", ctx)
	encryption, err := encrypt.NewSSEKMS(config.ServerConfigValues.Minio.EncryptionKeyID, ctx)
	if err != nil {
		log.Println(err)
		return minio.UploadInfo{}, err
	}

	// Upload the test file with FPutObject
	info, err := storage.MinioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ServerSideEncryption: encryption})
	if err != nil {
		log.Println(err)
		return minio.UploadInfo{}, err
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
	return info, nil
}

func CreateBucket(bucketName string) error {
	found, err := storage.MinioClient.BucketExists(context.Background(), bucketName)
	if err != nil {
		log.Println(err)
		return err
	}
	if found {
		log.Println("Bucket", bucketName, "already exists")
		return nil
	}

	// Create a bucket at region 'us-east-1' with object locking enabled.
	err = storage.MinioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{Region: "us-east-1", ObjectLocking: true})
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Successfully created bucket:", bucketName)
	return nil
}

func RemoveObject(object string, bucket string) error {
	opts := minio.RemoveObjectOptions{
		GovernanceBypass: true,
		VersionID:        "", // remove latest object version
	}

	err := storage.MinioClient.RemoveObject(context.Background(), bucket, object, opts)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Printf("Successfully removed object %s from bucket %s", object, bucket)
	return nil

}
