package dto

import "errors"

type UpdateFileRequest struct {
	BucketName  string `json:"bucketName"`
	Location    string `json:"location"`
	ObjectName  string `json:"objectName"`
	Filepath    string `json:"filepath"`
	ContentType string `json:"contentType"`
}

func (r *UpdateFileRequest) Validate() error {
	var errorMsg string
	if r.BucketName == "" {
		errorMsg = "Insert valid bucket name"
		err := errors.New(errorMsg)
		return err
	}

	if r.Location == "" {
		errorMsg = "Insert valid location"
		err := errors.New(errorMsg)
		return err
	}

	if r.ObjectName == "" {
		errorMsg = "Insert valid object name"
		err := errors.New(errorMsg)
		return err

	}

	if r.Filepath == "" {
		errorMsg = "Insert valid filepath"
		err := errors.New(errorMsg)
		return err

	}

	if r.ContentType == "" {
		errorMsg = "Insert valid content type"
		err := errors.New(errorMsg)
		return err

	}
	return nil
}
