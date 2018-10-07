package main

import (
	"context"
	"fmt"
	"image"
	"io"
	"log"
	"mime/multipart"
	"time"

	"cloud.google.com/go/storage"

	_ "image/jpeg"
	_ "image/png"

	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
	uuid "github.com/satori/go.uuid"
)

// FileName for a saved image.
func FileName(ext imaging.Format) (string, error) {
	uid := uuid.NewV4().String()
	d := time.Now().Format("2006-01-02-03-04-05")

	switch ext {
	case imaging.JPEG:
		return fmt.Sprintf("%s-%s.jpg", d, uid), nil
	case imaging.PNG:
		return fmt.Sprintf("%s-%s.png", d, uid), nil
	default:
		return "", fmt.Errorf("this format is not supported: %v", ext)
	}
}

// GetOrientation returns the orientation integer.
// See: https://github.com/disintegration/imaging/issues/30
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

// Manipulate the image based on config and ot.
func Manipulate(f io.Reader, config ImageConfig, ot int) (image.Image, error) {
	img, err := imaging.Decode(f)
	if err != nil {
		log.Printf("error when decoding the image for imaging lib: %s", err)
		return img, err
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
		img = imaging.Fill(
			img,
			config.Width,
			config.Height,
			imaging.Center,
			imaging.Lanczos,
		)
	} else {
		img = imaging.Fit(img, config.Width, config.Height, imaging.Lanczos)
	}

	return img, nil
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

// UploadFile a file to GCP.
func UploadFile(config ImageConfig, fh *multipart.FileHeader) (*ImageConfig, error) {
	mimeType := fh.Header.Get("Content-Type")
	formats := map[string]imaging.Format{
		"image/png":  imaging.PNG,
		"image/jpeg": imaging.JPEG,
		"image/jpg":  imaging.JPEG,
		"image/gif":  imaging.GIF,
	}
	format := formats[mimeType]

	name, err := FileName(format)
	if err != nil {
		return nil, err
	}

	publicURL := fmt.Sprintf(
		"https://storage.googleapis.com/%s/%s",
		GetEnv("BUCKET_ID", "default"),
		name,
	)

	f, err := fh.Open()
	if err != nil {
		return nil, err
	}

	defer f.Close()

	// Doesn't really matter if this fails. Doing our best to orientate.
	ot, _ := GetOrientation(f)

	seekBack(&f)

	img, err := Manipulate(f, config, ot)
	if err != nil {
		return nil, err
	}

	obj := App.bh.Object(name)

	w := obj.NewWriter(context.Background())
	w.ACL = []storage.ACLRule{{
		Entity: storage.AllUsers,
		Role:   storage.RoleReader,
	}}
	w.ContentType = mimeType
	w.CacheControl = fmt.Sprintf(
		"public, max-age=%s",
		GetEnv("CACHE_MAX_AGE", "86400"),
	)

	if err := imaging.Encode(w, img, format); err != nil {
		log.Println("error when writing to cloud")
		w.Close()
		return nil, err
	}

	// Close the connection.
	if err := w.Close(); err != nil {
		log.Println("error when closing connection")
		return nil, err
	}

	return &ImageConfig{
		FileName: name,
		Name:     config.Name,
		Width:    config.Width,
		Height:   config.Height,
		URL:      publicURL,
	}, nil
}
