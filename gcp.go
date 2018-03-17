package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"

	"cloud.google.com/go/storage"
	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
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

// Return the orientation integer.
// This is inspired from: https://github.com/disintegration/imaging/issues/30
func getOrientation(f io.Reader) (ot int, err error) {
	defer func() {
		if r := recover(); r != nil {
			ot = 1
		}
	}()

	x, err := exif.Decode(f)
	if err != nil {
		return 1, nil
	}

	tag, err := x.Get(exif.Orientation)
	if err != nil {
		return 1, err
	}

	ot, _ = tag.Int(0)

	return int(ot), nil
}

// Grab the width and height of an image in the cloud.
func getDimensions(name string) (int, int, error) {
	ctx := context.Background()
	rc, err := App.bh.Object(name).NewReader(ctx)

	if err != nil {
		return 0, 0, err
	}

	defer rc.Close()

	img, err := imaging.Decode(rc)

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
		log.Println("error when decoding")
		return ImageConfig{}, err
	}

	// Prep the writer for the gcp object.
	ctx := context.Background()
	w := App.bh.Object(name).NewWriter(ctx)
	w.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}
	w.ContentType = "image/jpeg"
	w.CacheControl = fmt.Sprintf("public, max-age=%s", GetEnv("CACHE_MAX_AGE", "86400"))

	// Get the dimensions from the exif data and rotate accordingly.
	f2, err := fh.Open()
	if err != nil {
		return ImageConfig{}, err
	}

	ot, err := getOrientation(f2)
	if err != nil {
		return ImageConfig{}, err
	}

	switch {
	case ot == 2:
		img = imaging.FlipH(img)
	case ot == 3:
		img = imaging.Rotate180(img)
	case ot == 4:
		img = imaging.FlipH(img)
		img = imaging.Rotate180(img)
	case ot == 5:
		img = imaging.FlipV(img)
		img = imaging.Rotate270(img)
	case ot == 6:
		img = imaging.Rotate270(img)
	case ot == 7:
		img = imaging.FlipV(img)
		img = imaging.Rotate90(img)
	case ot == 8:
		img = imaging.Rotate90(img)
	}

	// Do image processing while at the same time writing to the cloud.
	if config.Fill == true {
		img = imaging.Fill(img, config.Width, config.Height, imaging.Center, imaging.Lanczos)
	} else {
		img = imaging.Fit(img, config.Width, config.Height, imaging.Lanczos)
	}

	if err := imaging.Encode(w, img, imaging.JPEG); err != nil {
		log.Println("error when writing to cloud")
		return ImageConfig{}, err
	}

	// Close the connection.
	if err := w.Close(); err != nil {
		log.Println("error when closing connection")
		return ImageConfig{}, err
	}

	width, height, err := getDimensions(name)

	if err != nil {
		go RemoveFile(name)
		return ImageConfig{}, err
	}

	publicURL := fmt.Sprintf(
		"https://storage.googleapis.com/%s/%s",
		GetEnv("BUCKET_ID", "default"),
		name,
	)

	return ImageConfig{
		FileName: name,
		Name:     config.Name,
		Width:    width,
		Height:   height,
		URL:      publicURL,
	}, nil
}
