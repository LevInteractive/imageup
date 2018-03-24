package main

import (
	"io"
	"os"
	"testing"
)

func openFile(filename string) io.Reader {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	return f
}

func TestGetOrientation(t *testing.T) {
	tests := []struct {
		filename string
		expected int
	}{
		{"./testdata/ben.jpg", 1},
		{"./testdata/acai.jpg", 6},
		{"./testdata/dog.jpg", 6},
	}

	for _, tt := range tests {
		result, err := GetOrientation(openFile(tt.filename))
		if err != nil {
			t.Errorf("Received unexpected error when getting orientation: %s", err)
		}
		if result != tt.expected {
			t.Errorf("Expected GetOrientation(%s) to be %d", tt.filename, tt.expected)
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
