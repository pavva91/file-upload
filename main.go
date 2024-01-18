package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/pavva91/file-upload/api"
	"github.com/pavva91/file-upload/config"
	"github.com/pavva91/file-upload/storage"
)

func main() {
	env := os.Getenv("SERVER_ENVIRONMENT")

	log.Println(fmt.Sprintf("Running Environment: %s", env))

	switch env {
	case "dev":
		setConfig("./config/dev-config.yml")
	case "stage":
		log.Panic(fmt.Sprintf("Incorrect Dev Environment: %s\nInterrupt execution", env))
	case "prod":
		log.Panic(fmt.Sprintf("Incorrect Dev Environment: %s\nInterrupt execution", env))
	default:
		log.Panic(fmt.Sprintf("Incorrect Dev Environment: %s\nInterrupt execution", env))
	}

	storage.MinioClient = storage.CreateMinioClient()

	// Create a new request multiplexer
	// Take incoming requests and dispatch them to the matching handlers
	mux := http.NewServeMux()

	// Register the routes and handlers
	mux.Handle("/", &homeHandler{})
	mux.Handle("/health", &healthHandler{})
	mux.Handle("/files", &api.FilesHandler{})
	mux.Handle("/files/", &api.FilesHandler{})

	// Run the server
	fmt.Printf("Server is running on port %s", config.ServerConfigValues.Server.Port)
	// http.ListenAndServe(":8080", mux)
	http.ListenAndServe(fmt.Sprintf(":%s", config.ServerConfigValues.Server.Port), mux)
}

func setConfig(path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&config.ServerConfigValues)
	if err != nil {
		log.Fatal(err)
	}
}

type homeHandler struct{}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("home"))
}

type healthHandler struct{}

func (h *healthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("healthy"))
}
