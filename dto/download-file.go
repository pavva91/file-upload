package dto

import "errors"

type DownloadFileRequest struct {
	BucketName   string `json:"bucketName"`
	DownloadPath string `json:"downloadPath"`
}

func (r *DownloadFileRequest) Validate() error {
	var errorMsg string

	if r.BucketName == "" {
		errorMsg = "Insert valid bucket name"
		err := errors.New(errorMsg)
		return err
	}

	if r.DownloadPath == "" {
		errorMsg = "Insert valid download folder path"
		err := errors.New(errorMsg)
		return err

	}

	return nil
}
