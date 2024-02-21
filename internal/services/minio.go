// FileUploader.go MinIO example
package services

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/cheggaaa/pb"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/encrypt"
	"github.com/pavva91/file-upload/config"
	"github.com/pavva91/file-upload/internal/storage"
)

func EncryptAndUploadFileMultipart(objectName string, filePath string, contentType string, bucketName string) (minio.UploadInfo, error) {
	ctx := context.Background()

	// encryption, err := encrypt.NewSSEKMS("dev-key2", ctx)
	encryption, err := encrypt.NewSSEKMS(config.ServerConfigValues.Minio.EncryptionKeyID, ctx)
	if err != nil {
		log.Println(err)
		return minio.UploadInfo{}, err
	}

	sizeMiB := uint64(config.ServerConfigValues.Minio.FileChunkSize)

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

	// Progress reader is notified as PutObject makes progress with
	// the Reads inside.
	progress := pb.New64(fileStat.Size())
	progress.Start()

	isMultipart := false

	go func(bool) {
		notificationChan := storage.MinioClient.ListenNotification(ctx, "", "", []string{
			"s3:ObjectCreated:CompleteMultipartUpload",
		})
		notification := <-notificationChan
		log.Println(notification)
		isMultipart = true
	}(isMultipart)

	opts := minio.PutObjectOptions{
		DisableMultipart:     !config.ServerConfigValues.Minio.EnableMultipartUpload,
		PartSize:             1024 * 1024 * sizeMiB,
		ServerSideEncryption: encryption,
		Progress:             progress,
	}

	uploadInfo, err := storage.MinioClient.PutObject(ctx, bucketName, objectName, file, fileStat.Size(), opts)
	if err != nil {
		log.Println(err)
		return minio.UploadInfo{}, err
	}

	log.Printf("Successfully uploaded %s of size %d Bytes\n", objectName, uploadInfo.Size)

	time.Sleep(100 * time.Millisecond)
	if isMultipart {
		log.Println("used multipart upload for object", objectName)
	} else {
		log.Println("did not use multipart upload for object", objectName)
	}

	return uploadInfo, nil
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
	region := config.ServerConfigValues.Minio.Region
	err = storage.MinioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{Region: region, ObjectLocking: true})
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
