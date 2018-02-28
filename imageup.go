package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/storage"
)

type errorJSON struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// FileUpload represents both the input and output for a file.
type FileUpload struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Fill   bool   `json:"fill"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// AppConfig for the global singleton
type AppConfig struct {
	bh *storage.BucketHandle
}

// App is the singleton.
var App = &AppConfig{}

// Grab a storage bucket client instance. This will persist throughout the
// lifetime of the server.
func configureStorage(bucketID string) (*storage.BucketHandle, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.Bucket(bucketID), nil
}

// GetEnv grabs env with a fallback
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Provision the response for a json payload.
func jsonResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

// UploadFile uploads a file to the server
func handleRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", GetEnv("CORS", "*"))

	// Validate the incoming request.
	if r.Method != http.MethodPost {
		log.Printf("Error with the request (non-POST)")
		jsonResponse(w, http.StatusMethodNotAllowed, errorJSON{
			http.StatusMethodNotAllowed,
			"fail",
		})
		return
	}

	// Grab the file from the request.
	file, handle, err := r.FormFile("file")
	if err != nil {
		log.Printf("Error retrieving file: %v", err)
		jsonResponse(w, http.StatusBadRequest, errorJSON{
			http.StatusBadRequest,
			"problem finding the file",
		})
		return
	}

	defer file.Close()

	// Parse the file configuration json array.
	var fileConfig []FileUpload

	if json.Unmarshal([]byte(r.FormValue("sizes")), &fileConfig) != nil {
		log.Printf("Error handling params: %v", err)
		log.Printf("This is what the config looks like: %v", r.FormValue("sizes"))
		jsonResponse(w, http.StatusNotAcceptable, errorJSON{
			http.StatusBadRequest,
			"there is a problem with the size configuration",
		})
		return
	}
	log.Printf("the files: %v", fileConfig)

	if len(fileConfig) < 1 {
		log.Printf("No size sent with request.")
		jsonResponse(w, http.StatusNotAcceptable, errorJSON{
			http.StatusBadRequest,
			"there were no size instructions sent with request",
		})
		return
	}

	log.Printf("These are the sizes: %v", fileConfig)

	var uploadedFiles []FileUpload

	// Handle the file uploads.
	for fileConfig := range conf {
		path, err := UploadFile(file, conf, handle)
		if err != nil {
			log.Printf("Error uploading file: %v", err)
			jsonResponse(w, http.StatusBadRequest, errorJSON{
				http.StatusBadRequest,
				"invalid",
			})
			return
		}
		uploadedFiles = append(uploadedFiles, path)
	}

	log.Printf("Value: %s", r.FormValue("sizes"))
	log.Printf("path: %s", path)
	jsonResponse(w, http.StatusCreated, &FileUpload{})
}

func main() {
	port := GetEnv("SERVER_PORT", "31111")
	bucket := GetEnv("BUCKET_ID", "worktaps-dev")

	bh, err := configureStorage(bucket)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on port %s. Using bucket %s.", port, bucket)
	app.bh = bh

	http.HandleFunc("/", handleRequest)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
