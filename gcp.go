package main

import (
	"context"
	"fmt"
	"image"
	"io"
	"log"
	"mime/multipart"

	_ "image/jpeg"

	"cloud.google.com/go/storage"
	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
	uuid "github.com/satori/go.uuid"
)

// GetOrientation returns the orientation integer.
// This is inspired from: https://github.com/disintegration/imaging/issues/30
func GetOrientation(f io.Reader) (ot int, err error) {
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

// GetDimensions grabs the width and height of an image in the cloud.
func GetDimensions(f io.Reader) (int, int, error) {
	img, _, err := image.Decode(f)
	if err != nil {
		return 0, 0, err
	}

	bounds := img.Bounds()
	return bounds.Max.X, bounds.Max.Y, nil
}

// RemoveFileFromGCP removes a file from GCP.
func RemoveFileFromGCP(name string) error {
	log.Printf("Requesting to delete %s", name)
	ctx := context.Background()
	err := App.bh.Object(name).Delete(ctx)

	if err != nil {
		log.Printf("There was an error removing image from storage: %v", err)
	}

	return err
}

// Seek back to the beginning of the file. This should be called after each
// read.
func seekBack(f *multipart.File) error {
	if _, err := (*f).Seek(0, 0); err != nil {
		log.Printf("error seeking back: %v", err)
		return err
	}
	return nil
}

// UploadFile adds a file to GCP.
func UploadFile(config ImageConfig, fh *multipart.FileHeader) (ImageConfig, error) {
	ctx := context.Background()

	// Grab the file reader.
	f, err := fh.Open()
	if err != nil {
		return ImageConfig{}, err
	}

	defer f.Close()

	// This will be what the image is saved as. Using the uuid to create a unique
	// ID. This is good enough for us for now.
	name := fmt.Sprintf("%s-%s.jpg", config.Name, uuid.NewV4().String())

	// We know what the public url is going to be before even uploading it.
	publicURL := fmt.Sprintf(
		"https://storage.googleapis.com/%s/%s",
		GetEnv("BUCKET_ID", "default"),
		name,
	)

	img, err := imaging.Decode(f)
	if err != nil {
		log.Printf("error when decoding the image for imaging lib: %s", err)
		return ImageConfig{}, err
	}

	// Set the seeker back so we can continue to read from the file.
	seekBack(&f)

	// Get the orientation of the image and rotate accordingly.
	ot, err := GetOrientation(f)
	if err != nil {
		log.Printf("error when getting orientation from exif info: %s", err)
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

	obj := App.bh.Object(name)

	// Prep the writer for the gcp object.
	w := obj.NewWriter(ctx)
	w.ACL = []storage.ACLRule{{
		Entity: storage.AllUsers,
		Role:   storage.RoleReader,
	}}
	w.ContentType = "image/jpeg"
	w.CacheControl = fmt.Sprintf("public, max-age=%s", GetEnv("CACHE_MAX_AGE", "86400"))

	if err := imaging.Encode(w, img, imaging.JPEG); err != nil {
		log.Println("error when writing to cloud")
		w.Close()
		return ImageConfig{}, err
	}

	// Close the connection.
	if err := w.Close(); err != nil {
		log.Println("error when closing connection")
		return ImageConfig{}, err
	}

	return ImageConfig{
		FileName: name,
		Name:     config.Name,
		Width:    config.Width,
		Height:   config.Height,
		URL:      publicURL,
	}, nil
}
