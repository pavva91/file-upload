package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/pavva91/file-upload/dto"
	"github.com/pavva91/file-upload/errorhandlers"
	"github.com/pavva91/file-upload/services"
	"github.com/pavva91/file-upload/storage"
)

type FilesHandler struct{}

var (
	FileRe         = regexp.MustCompile(`^/files/*$`)
	FileReWithID   = regexp.MustCompile(`^/files/([a-z0-9]+(?:-[a-z0-9]+)+)$`)
	FileReWithName = regexp.MustCompile(`^/files/.+$`)
)

func (h *FilesHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	var reqBody dto.UploadFileRequest

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = reqBody.Validate()
	if err != nil {
		errorhandlers.BadRequestHandler(w, r, err)
		return
	}

	bucketName := reqBody.BucketName

	bucketExists, err := services.BucketExist(bucketName)
	if err != nil {
		log.Println(err.Error())
		errorhandlers.InternalServerErrorHandler(w, r)
		return
	}

	if !bucketExists {
		msg := fmt.Sprintln("bucket", bucketName, "does not exist")
		err := errors.New(msg)
		log.Println(err.Error())
		errorhandlers.BadRequestHandler(w, r, err)
		return
	}

	// Upload the test file
	// Change the value of filePath if the file is in another location

	objectName := reqBody.ObjectName
	filePath := reqBody.Filepath
	contentType := reqBody.ContentType

	// services.UploadFile(objectName, filePath, contentType, bucketName)
	services.EncryptAndUploadFile(objectName, filePath, contentType, bucketName)
	// services.EncryptAndUploadFileMultipart(objectName, filePath, contentType, bucketName)
	if err != nil {
		log.Println(err)
		errorhandlers.InternalServerErrorHandler(w, r)
		return
	}

	// w.WriteHeader(http.StatusOK)
}

func (h *FilesHandler) DownloadFile(w http.ResponseWriter, r *http.Request) {
	var reqBody dto.DownloadFileRequest

	fileName := strings.TrimPrefix(r.URL.Path, "/files/")
	log.Println(fmt.Sprintf("Request download file: %s", fileName))

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		if err.Error() == "EOF" {
			err = errors.New("No Request JSON Body")
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = reqBody.Validate()
	if err != nil {
		errorhandlers.BadRequestHandler(w, r, err)
		return
	}

	bucket := reqBody.BucketName

	bucketExists, err := services.BucketExist(bucket)
	if err != nil {
		log.Println(err.Error())
		errorhandlers.InternalServerErrorHandler(w, r)
		return
	}

	if !bucketExists {
		msg := fmt.Sprintln("bucket", bucket, "does not exist")
		err := errors.New(msg)
		log.Println(err.Error())
		errorhandlers.BadRequestHandler(w, r, err)
		return
	}

	downloadPath := reqBody.DownloadPath

	err = services.DownloadFile(bucket, fileName, downloadPath)
	if err != nil {
		log.Println(err)
		if err.Error() == "The specified key does not exist." {
			err = errors.New(fmt.Sprintf("Specified file %s is not present in bucket %s", fileName, bucket))
			log.Println(err)
			errorhandlers.BadRequestHandler(w, r, err)
		} else {
			errorhandlers.InternalServerErrorHandler(w, r)
		}
		return
	}

	msg := fmt.Sprintf("File %s correctly downloaded in: %s", fileName, downloadPath)
	log.Println(msg)
	w.Write([]byte(msg))
}

func (h *FilesHandler) ListFiles(w http.ResponseWriter, r *http.Request) {
	var objects []minio.ObjectInfo
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	objectCh := storage.MinioClient.ListObjects(ctx, "testbucket", minio.ListObjectsOptions{
		Prefix:    "",
		Recursive: true,
	})
	for o := range objectCh {
		if o.Err != nil {
			fmt.Println(o.Err)
			return
		}
		fmt.Println(o)
		objects = append(objects, o)
	}

	// profile := Profile{"Alex", []string{"snowboarding", "programming"}}

	js, err := json.Marshal(objects)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

	w.WriteHeader(http.StatusOK)
}

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
	case r.Method == http.MethodGet && FileReWithName.MatchString(r.URL.Path):
		h.DownloadFile(w, r)
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
