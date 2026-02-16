package services

import (
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// createTestJPEG writes a small JPEG to path and returns the path.
func createTestJPEG(t *testing.T, dir, name string) string {
	t.Helper()
	path := filepath.Join(dir, "images", name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	img := image.NewRGBA(image.Rect(0, 0, 20, 20))
	for x := 0; x < 20; x++ {
		for y := 0; y < 20; y++ {
			img.Set(x, y, color.RGBA{R: 100, G: 150, B: 200, A: 255})
		}
	}
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := jpeg.Encode(f, img, nil); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestRunPipeline(t *testing.T) {
	dir := t.TempDir()
	p1 := createTestJPEG(t, dir, "test1.jpg")
	p2 := createTestJPEG(t, dir, "test2.jpg")

	svc := NewService()
	svc.Run([]string{p1, p2})

	// The pipeline replaces "images/" → "images/output/" in the path.
	for _, name := range []string{"test1.jpg", "test2.jpg"} {
		inputPath := filepath.Join(dir, "images", name)
		outPath := strings.Replace(inputPath, "images/", "images/output/", 1)
		if _, err := os.Stat(outPath); err != nil {
			t.Errorf("expected output file %s: %v", outPath, err)
		}
	}
}

func TestRunPipelineInvalidPath(t *testing.T) {
	svc := NewService()
	// Should not panic; errors are handled gracefully.
	svc.Run([]string{"/nonexistent/image.jpg"})
}
