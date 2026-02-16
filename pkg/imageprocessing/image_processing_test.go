package imageprocessing

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"
)

// newTestImage creates a small RGBA image for testing.
func newTestImage(w, h int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			img.Set(x, y, color.RGBA{R: uint8(x % 256), G: uint8(y % 256), B: 128, A: 255})
		}
	}
	return img
}

func TestGrayscale(t *testing.T) {
	src := newTestImage(10, 10)
	gray := Grayscale(src)

	if gray.Bounds() != src.Bounds() {
		t.Fatalf("bounds mismatch: got %v, want %v", gray.Bounds(), src.Bounds())
	}

	// Every pixel should be a shade of gray (R == G == B).
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			r, g, b, _ := gray.At(x, y).RGBA()
			if r != g || g != b {
				t.Fatalf("pixel (%d,%d) is not gray: r=%d g=%d b=%d", x, y, r, g, b)
			}
		}
	}
}

func TestResize(t *testing.T) {
	src := newTestImage(100, 80)
	resized := Resize(src, 50, 40)

	bounds := resized.Bounds()
	if bounds.Dx() != 50 || bounds.Dy() != 40 {
		t.Fatalf("unexpected size: %dx%d, want 50x40", bounds.Dx(), bounds.Dy())
	}
}

func TestWriteAndReadImage(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.jpg")

	src := newTestImage(20, 20)
	if err := WriteImage(path, src); err != nil {
		t.Fatalf("WriteImage: %v", err)
	}

	img, err := ReadImage(path)
	if err != nil {
		t.Fatalf("ReadImage: %v", err)
	}

	if img.Bounds().Dx() != 20 || img.Bounds().Dy() != 20 {
		t.Fatalf("unexpected size: %v", img.Bounds())
	}
}

func TestWriteImageCreatesDirectories(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "a", "b", "test.jpg")

	src := newTestImage(5, 5)
	if err := WriteImage(path, src); err != nil {
		t.Fatalf("WriteImage with nested dirs: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("file not created: %v", err)
	}
}

func TestReadImageNotFound(t *testing.T) {
	_, err := ReadImage("/nonexistent/path.jpg")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
