// Package imageprocessing provides functions for reading, writing, and
// transforming images. It supports JPEG decoding/encoding, resizing via
// Lanczos3 interpolation, and grayscale conversion.
package imageprocessing

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
)

// ReadImage decodes a JPEG image from the file at path and returns it.
// Returns an error if the file cannot be opened or decoded.
func ReadImage(path string) (image.Image, error) {
	inputFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer inputFile.Close()

	img, _, err := image.Decode(inputFile)
	if err != nil {
		return nil, fmt.Errorf("decode %s: %w", path, err)
	}

	return img, nil
}

// WriteImage encodes img as JPEG and writes it to the file at path.
// Parent directories are created automatically if they do not exist.
func WriteImage(path string, img image.Image) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", filepath.Dir(path), err)
	}

	outputFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create %s: %w", path, err)
	}
	defer outputFile.Close()

	if err := jpeg.Encode(outputFile, img, nil); err != nil {
		return fmt.Errorf("encode %s: %w", path, err)
	}

	return nil
}

// Grayscale converts img to an 8-bit grayscale image using the standard
// luminance formula.
func Grayscale(img image.Image) image.Image {
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			grayImg.Set(x, y, color.GrayModel.Convert(img.At(x, y)))
		}
	}

	return grayImg
}

// Resize scales img to the given width and height using Lanczos3 interpolation.
func Resize(img image.Image, width, height uint) image.Image {
	return resize.Resize(width, height, img, resize.Lanczos3)
}
