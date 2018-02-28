package main

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path"

	"cloud.google.com/go/storage"
	uuid "github.com/satori/go.uuid"
)

// RemoveFile removes a file from GCP.
func RemoveFile() {
	// todo
}

// UploadFile adds a file to GCP.
func UploadFile(f multipart.File, config []FileUpload, fh *multipart.FileHeader) (FileUpload, error) {
	ctx := context.Background()
	name := uuid.NewV4().String() + path.Ext(fh.Filename)
	w := App.bh.Object(name).NewWriter(ctx)
	w.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}
	w.ContentType = fh.Header.Get("Content-Type")
	w.CacheControl = fmt.Sprintf("public, max-age=%s", GetEnv("CACHE_MAX_AGE"))

	if _, err := io.Copy(w, f); err != nil {
		return "", err
	}

	if err := w.Close(); err != nil {
		return "", err
	}

	publicURL := fmt.Sprintf(
		"https://storage.googleapis.com/%s/%s",
		"worktaps-dev",
		name,
	)

	return FileUpload{
		Name:   config.Name,
		Width:  config.Width,
		Height: config.Height,
		URL:    publicURL,
	}, nil
}
