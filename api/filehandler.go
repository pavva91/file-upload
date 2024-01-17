package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/minio/minio-go/v7"
	"github.com/pavva91/file-upload/dto"
	"github.com/pavva91/file-upload/errorhandlers"
	"github.com/pavva91/file-upload/storage"
)

type FilesHandler struct{}

var (
	FileRe                     = regexp.MustCompile(`^/files/*$`)
	FileReWithID               = regexp.MustCompile(`^/files/([a-z0-9]+(?:-[a-z0-9]+)+)$`)
)

func (h *FilesHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var req dto.UpdateFileRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = req.Validate()
	if (err != nil) {
		errorhandlers.BadRequestHandler(w, r, err)
		return
	}

	// Make a new bucket called testbucket.

	// bucketName := "testbucket"
	// location := "us-east-1"

	bucketName := req.BucketName
	location := req.Location

	// Upload the test file
	// Change the value of filePath if the file is in another location

	// objectName := "testdata"
	// filePath := "/tmp/testdata"
	// contentType := "application/octet-stream"

	objectName := req.ObjectName
	filePath := req.Filepath
	contentType := req.ContentType

	err = storage.MinioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := storage.MinioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			msg := fmt.Sprintf("We already own %s\n", bucketName)
			log.Printf(msg)
			// errorhandlers.BadRequestHandler(w, r, msg)
		} else {
			log.Fatalln(err)
			errorhandlers.InternalServerErrorHandler(w, r)
			return
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}

	// Upload the test file with FPutObject
	info, err := storage.MinioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
		errorhandlers.BadRequestHandler(w, r, err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
	w.WriteHeader(http.StatusOK)
}

func (h *FilesHandler) ListFiles(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	objectCh := storage.MinioClient.ListObjects(ctx, "testbucket", minio.ListObjectsOptions{
		Prefix:    "testdata",
		Recursive: true,
	})
	for object := range objectCh {
		if object.Err != nil {
			fmt.Println(object.Err)
			return
		}
		fmt.Println(object)
	}

	// buckets, err := minioClient.ListBuckets(context.Background())
	// if err != nil {
	// 	log.Fatalln(err)
	// 	errorhandlers.InternalServerErrorHandler(w, r)
	// }
	// for _, bucket := range buckets {
	// 	fmt.Println(bucket)
	// }

	w.WriteHeader(http.StatusOK)
	// w.Write(buckets)
}
func (h *FilesHandler) GetFile(w http.ResponseWriter, r *http.Request)    {}
func (h *FilesHandler) UpdateFile(w http.ResponseWriter, r *http.Request) {}
func (h *FilesHandler) DeleteFile(w http.ResponseWriter, r *http.Request) {}

func (h *FilesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && FileRe.MatchString(r.URL.Path):
		h.UploadFile(w, r)
		return
	case r.Method == http.MethodGet && FileRe.MatchString(r.URL.Path):
		h.ListFiles(w, r)
		return
	case r.Method == http.MethodGet && FileReWithID.MatchString(r.URL.Path):
		h.GetFile(w, r)
		return
	case r.Method == http.MethodPut && FileReWithID.MatchString(r.URL.Path):
		h.UpdateFile(w, r)
		return
	case r.Method == http.MethodDelete && FileReWithID.MatchString(r.URL.Path):
		h.DeleteFile(w, r)
		return
	default:
		return
	}
}
