package main

import (
	"context"
	"fmt"
	"image/jpeg"
	"log"
	"mime/multipart"

	"cloud.google.com/go/storage"
	"github.com/disintegration/imaging"
	uuid "github.com/satori/go.uuid"
)

// RemoveFile removes a file from GCP.
func RemoveFile(name string) error {
	log.Printf("Requesting to delete %s", name)
	ctx := context.Background()
	err := App.bh.Object(name).Delete(ctx)

	if err != nil {
		log.Printf("There was an error removing image from storage: %v", err)
	}

	return err
}

// Grab the width and height of an image in the cloud.
func getDimensions(name string) (int, int, error) {
	ctx := context.Background()
	rc, err := App.bh.Object(name).NewReader(ctx)

	if err != nil {
		return 0, 0, err
	}

	defer rc.Close()

	img, err := jpeg.Decode(rc)

	if err != nil {
		return 0, 0, err
	}

	bounds := img.Bounds()

	return bounds.Max.X, bounds.Max.Y, nil
}

// UploadFile adds a file to GCP.
func UploadFile(f multipart.File, config ImageConfig, fh *multipart.FileHeader) (ImageConfig, error) {
	name := fmt.Sprintf("%s-%s.jpg", config.Name, uuid.NewV4().String())
	img, err := imaging.Decode(f)

	if err != nil {
		return ImageConfig{}, err
	}

	// Prep the writer for the gcp object.
	ctx := context.Background()
	w := App.bh.Object(name).NewWriter(ctx)
	w.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}
	w.ContentType = "image/jpeg"
	w.CacheControl = fmt.Sprintf("public, max-age=%s", GetEnv("CACHE_MAX_AGE", "86400"))

	// Do image processing while at the same time writing to the cloud.
	if config.Fill == true {
		img = imaging.Fill(img, config.Width, config.Height, imaging.Center, imaging.Lanczos)
	} else {
		img = imaging.Fit(img, config.Width, config.Height, imaging.Lanczos)
	}

	if err := imaging.Encode(w, img, imaging.JPEG); err != nil {
		return ImageConfig{}, err
	}

	// Close the connection.
	if err := w.Close(); err != nil {
		return ImageConfig{}, err
	}

	width, height, err := getDimensions(name)

	if err != nil {
		go RemoveFile(name)
		return ImageConfig{}, err
	}

	publicURL := fmt.Sprintf(
		"https://storage.googleapis.com/%s/%s",
		"worktaps-dev",
		name,
	)

	return ImageConfig{
		Name:   config.Name,
		Width:  width,
		Height: height,
		URL:    publicURL,
	}, nil
}
