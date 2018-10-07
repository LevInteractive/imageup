package main

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/disintegration/imaging"
)

func openFile(filename string) io.Reader {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	return f
}

func TestFileName(t *testing.T) {
	tests := []struct {
		format imaging.Format
		needle string
	}{
		{imaging.JPEG, "jpg"},
		{imaging.PNG, "png"},
	}

	for _, tt := range tests {
		f := FileName(tt.format)
		if strings.Contains(f, tt.needle) != true {
			t.Errorf("file (%s) did not contain format: %s", f, tt.needle)
		}
	}
}

func TestGetOrientation(t *testing.T) {
	tests := []struct {
		filename string
		expected int
	}{
		{"./testdata/ben.jpg", 1},
		{"./testdata/acai.jpg", 6},
		{"./testdata/dog.jpg", 6},
		{"./testdata/small-png.png", 1},
		{"./testdata/small-logo.jpg", 1},
	}

	for _, tt := range tests {
		result, err := GetOrientation(openFile(tt.filename))
		if err != nil {
			t.Errorf("Received unexpected error when getting orientation: %s", err)
		}
		if result != tt.expected {
			t.Errorf("Expected GetOrientation(%s) to be %d but got %v", tt.filename, tt.expected, result)
		}
	}
}

func TestGetDimensions(t *testing.T) {
	tests := []struct {
		filename string
		width    int
		height   int
	}{
		{"./testdata/ben.jpg", 2000, 1500},
		{"./testdata/acai.jpg", 4032, 3024},
		{"./testdata/dog.jpg", 4592, 3448},
		{"./testdata/small-png.png", 167, 167},
	}

	for _, tt := range tests {
		width, height, err := GetDimensions(openFile(tt.filename))
		if err != nil {
			t.Errorf("Received unexpected error when getting dimensions: %s", err)
		}
		if width != tt.width {
			t.Errorf("GetDimensions(%s) has wrong width. Expected %d to be %d", tt.filename, width, tt.width)
		}
		if height != tt.height {
			t.Errorf("GetDimensions(%s) has wrong height dims. Expected %d to be %d", tt.filename, height, tt.height)
		}
	}
}

func TestManipulate(t *testing.T) {
	tests := []struct {
		filename string
		configs  ImageConfig
		ot       int
	}{
		{
			"./testdata/ben.jpg",
			ImageConfig{
				Fill:   false,
				Width:  100,
				Height: 100,
			},
			1,
		},
		{
			"./testdata/small-png.png",
			ImageConfig{
				Fill:   false,
				Width:  100,
				Height: 100,
			},
			4,
		},
	}

	for _, tt := range tests {
		_, err := Manipulate(openFile(tt.filename), tt.configs, 1)
		if err != nil {
			t.Errorf("received an error when manipulating image: %v", err)
		}
	}
}
