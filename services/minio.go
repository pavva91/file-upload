// FileUploader.go MinIO example
package services

import (
	"context"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/encrypt"
	"github.com/pavva91/file-upload/config"
	"github.com/pavva91/file-upload/storage"
)

func UploadFile(objectName string, filePath string, contentType string, bucketName string) (minio.UploadInfo, error) {
	ctx := context.Background()

	// Upload the test file
	// Change the value of filePath if the file is in another location

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

func EncryptAndUploadFileMultipart(objectName string, filePath string, contentType string, bucketName string) (minio.UploadInfo, error) {
	ctx := context.Background()

	// encryption, err := encrypt.NewSSEKMS("dev-key2", ctx)
	encryption, err := encrypt.NewSSEKMS(config.ServerConfigValues.Minio.EncryptionKeyID, ctx)
	if err != nil {
		log.Println(err)
		return minio.UploadInfo{}, err
	}

	sizeMiB := uint64(config.ServerConfigValues.Minio.FileChunkSize)

	// INFO: For multipart set minio.PutObjectOptions.DisableMultipart
	opts := minio.PutObjectOptions{
		DisableMultipart: !config.ServerConfigValues.Minio.EnableMultipartUpload,
		PartSize:         1024 * 1024 * sizeMiB,

		// NOTE: absMinPartSize - absolute minimum part size (5 MiB) below which
		// a part in a multipart upload may not be uploaded.
		// const absMinPartSize = 1024 * 1024 * 5
		// https://github.com/minio/minio-go/blob/6ad2b4a17816b1e991f73e598885c07704aea7ef/constants.go#L24

		// NOTE: DEFAULT Min Part Size if not defined
		// https://github.com/minio/minio-go/blob/6ad2b4a17816b1e991f73e598885c07704aea7ef/api-put-object.go#L302
		// minPartSize - minimum part size 16MiB per object after which
		// putObject behaves internally as multipart.
		// const minPartSize = 1024 * 1024 * 16 // 16MiB

		// NOTE: maxMultipartPutObjectSize - maximum size 5TiB of object for
		// Multipart operation.
		// const maxMultipartPutObjectSize = 1024 * 1024 * 1024 * 1024 * 5

		ServerSideEncryption: encryption,
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
		return minio.UploadInfo{}, err
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		log.Println(err)
		return minio.UploadInfo{}, err
	}
	
	// Upload the test file with FPutObject
	// uploadInfo, err := storage.MinioClient.FPutObject(ctx, bucketName, objectName, filePath, opts)
	// uploadInfo, err := storage.MinioClient.PutObject(ctx, bucketName, objectName, file, -1, opts)
	uploadInfo, err := storage.MinioClient.PutObject(ctx, bucketName, objectName, file, fileStat.Size(), opts)
	if err != nil {
		log.Println(err)
		return minio.UploadInfo{}, err
	}

	log.Printf("Successfully uploaded %s of size %d Bytes\n", objectName, uploadInfo.Size)
	return uploadInfo, nil
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
