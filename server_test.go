package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pavva91/file-upload/api"
	"github.com/pavva91/file-upload/services"
	"github.com/pavva91/file-upload/storage"
)

func TestFileUpload(t *testing.T) {

	setConfig("./config/dev-config.yml")
	storage.MinioClient = storage.CreateMinioClient()

	bucketName := "test"
	err := services.CreateBucket(bucketName)
	if err != nil {
		t.Fatal(err)
	}

	objectName := "object1"
	err = services.RemoveObject(objectName, bucketName)
	if err != nil {
		t.Fatal(err)
	}

	ts := httptest.NewServer(&api.FilesHandler{})
	// defer ts.Close()

	fileHandlerURL := ts.URL + "/files"

	newreq := func(method, url string, body io.Reader) *http.Request {
		r, err := http.NewRequest(method, url, body)
		if err != nil {
			t.Fatal(err)
		}
		return r
	}

	tests := map[string]struct {
		request  *http.Request
		response string
		status   int
	}{
		"POST /files nil body": {
			request:  newreq("POST", fileHandlerURL, nil),
			response: "EOF",
			status:   400,
		},
		"POST /files empty body": {
			request:  newreq("POST", fileHandlerURL, strings.NewReader(`{}`)),
			response: "Insert valid bucket name",
			status:   400,
		},
		"POST /files without bucket name": {
			request:  newreq("POST", fileHandlerURL, strings.NewReader(`{"objectName":"objectname"}`)),
			response: "Insert valid bucket name",
			status:   400,
		},
		"POST /files without object name": {
			request:  newreq("POST", fileHandlerURL, strings.NewReader(`{"bucketName":"bucketname"}`)),
			response: "Insert valid object name",
			status:   400,
		},
		"POST /files without filepath": {
			request:  newreq("POST", fileHandlerURL, strings.NewReader(`{"bucketName":"bucketname", "objectName":"objectname"}`)),
			response: "Insert valid filepath",
			status:   400,
		},
		"POST /files without content-type": {
			request:  newreq("POST", fileHandlerURL, strings.NewReader(`{"bucketName":"bucketname", "objectName":"objectname", "filepath":"/tmp/test.txt"}`)),
			response: "Insert valid content type",
			status:   400,
		},
		"POST /files wrong bucketname": {
			request:  newreq("POST", fileHandlerURL, strings.NewReader(`{"bucketName":"wrongbucketname", "objectName":"objectname", "filepath":"/tmp/test.txt", "contentType":"application/octet-stream"}`)),
			response: "bucket wrongbucketname does not exist",
			status:   400,
		},
		"POST /files Upload OK": {
			request:  newreq("POST", fileHandlerURL, strings.NewReader(fmt.Sprintf(`{"bucketName":"%s", "objectName":"%s", "filepath":"/tmp/test.txt", "contentType":"application/octet-stream"}`, bucketName, objectName))),
			response: fmt.Sprintf(""),
			status:   200,
		},
	}

	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actualResponse, err := http.DefaultClient.Do(test.request)
			if err != nil {
				t.Fatal(err)
			}

			defer actualResponse.Body.Close()
			// check for expected response here
			b, err := io.ReadAll(actualResponse.Body)

			if err != nil {
				log.Fatalln(err)
			}

			if !strings.Contains(string(b), test.response) {
				t.Errorf("got %s, want %s", string(b), test.response)
			}

			if actualResponse.StatusCode != test.status {
				t.Errorf("got %d, want %d", actualResponse.StatusCode, test.status)
			}

		})
	}
}
