package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/storage"
)

type jsonResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ImageConfig represents both the input and output for a file.
type ImageConfig struct {
	FileName string `json:"fileName"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	Fill     bool   `json:"fill"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
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

// Remove all files in array. This is used to clean up something that went
// wrong.
func removeAll(files []ImageConfig) {
	for _, conf := range files {
		go RemoveFile(conf.FileName)
	}
}

// UploadFile uploads a file to the server
func handleRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", GetEnv("CORS", "*"))

	switch r.Method {
	case http.MethodDelete:
		for _, fname := range strings.Split(r.FormValue("files"), ",") {
			go RemoveFile(strings.TrimSpace(fname))
		}

		jsonResponse(w, http.StatusOK, jsonResp{
			http.StatusOK,
			"file(s) queued to be removed",
		})
	case http.MethodPost:
		// Grab the file from the request.
		file, handle, err := r.FormFile("file")
		if err != nil {
			log.Printf("Error retrieving file: %v", err)
			jsonResponse(w, http.StatusBadRequest, jsonResp{
				http.StatusBadRequest,
				"problem finding the file",
			})
			return
		}

		defer file.Close()

		// Parse the file configuration json array.
		var configs []ImageConfig

		if json.Unmarshal([]byte(r.FormValue("sizes")), &configs) != nil {
			log.Printf("Error handling params: %v", err)
			log.Printf("This is what the config looks like: %v", r.FormValue("sizes"))
			jsonResponse(w, http.StatusNotAcceptable, jsonResp{
				http.StatusBadRequest,
				"there is a problem with the size configuration",
			})
			return
		}

		if len(configs) < 1 {
			log.Printf("No size sent with request.")
			jsonResponse(w, http.StatusNotAcceptable, jsonResp{
				http.StatusBadRequest,
				"there were no size instructions sent with request",
			})
			return
		}

		var uploadedFiles []ImageConfig

		// Handle the file uploads.
		for _, conf := range configs {
			log.Printf("Processing image with size: %v", conf)

			if _, err = file.Seek(0, os.SEEK_SET); err != nil {
				log.Printf("Error seeking file: %v", err)
				removeAll(uploadedFiles)
				jsonResponse(w, http.StatusBadRequest, jsonResp{
					http.StatusBadRequest,
					"error seeking file",
				})
				return
			}

			c, err := UploadFile(file, conf, handle)
			if err != nil {
				log.Printf("Error uploading file: %v", err)
				removeAll(uploadedFiles)
				jsonResponse(w, http.StatusBadRequest, jsonResp{
					http.StatusBadRequest,
					"unknown error while uploading",
				})
				return
			}

			uploadedFiles = append(uploadedFiles, c)
		}

		jsonResponse(w, http.StatusCreated, uploadedFiles)

	default:
		log.Printf("Error with the request (non-POST or non-DELETE)")
		jsonResponse(w, http.StatusMethodNotAllowed, jsonResp{
			http.StatusMethodNotAllowed,
			"fail",
		})
	}
}

func main() {
	port := GetEnv("SERVER_PORT", "31111")
	bucket := GetEnv("BUCKET_ID", "default")

	bh, err := configureStorage(bucket)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on port %s. Using bucket %s.", port, bucket)
	App.bh = bh

	http.HandleFunc("/", handleRequest)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
