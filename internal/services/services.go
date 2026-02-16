// Package services orchestrates image-processing pipelines using Go's
// concurrency patterns. Each pipeline stage runs in its own goroutine and
// communicates via channels, forming a linear pipeline:
//
//	load → resize → grayscale → save
package services

import (
	"fmt"
	"image"
	"strings"

	"github.com/revanthvermareddy/image-processing/pkg/imageprocessing"
)

// DefaultWidth and DefaultHeight define the target resize dimensions.
const (
	DefaultWidth  uint = 500
	DefaultHeight uint = 500
)

// Service defines the interface for running an image-processing pipeline.
type Service interface {
	Run(imagePaths []string)
}

// ServiceImpl is the default implementation of Service.
type ServiceImpl struct{}

// NewService returns a new ServiceImpl.
func NewService() Service {
	return &ServiceImpl{}
}

// Run executes the full image-processing pipeline for the given image paths.
// Each image is loaded, resized, converted to grayscale, and saved to an
// output directory. All stages run concurrently via channels.
func (s *ServiceImpl) Run(imagePaths []string) {
	channel1 := loadImage(imagePaths)
	channel2 := resize(channel1)
	channel3 := convertToGrayScale(channel2)
	writeResults := saveImage(channel3)

	for result := range writeResults {
		if result.Err != nil {
			fmt.Printf("Error processing %s: %v\n", result.Path, result.Err)
		} else {
			fmt.Printf("Saved: %s\n", result.Path)
		}
	}
}

// Job carries an image and its file paths through the pipeline stages.
type Job struct {
	InputPath string
	Image     image.Image
	OutPath   string
	Err       error
}

// Result captures the outcome of saving a processed image.
type Result struct {
	Path string
	Err  error
}

func loadImage(paths []string) <-chan Job {
	out := make(chan Job)
	go func() {
		for _, path := range paths {
			job := Job{
				InputPath: path,
				OutPath:   strings.Replace(path, "images/", "images/output/", 1),
			}
			img, err := imageprocessing.ReadImage(path)
			if err != nil {
				job.Err = fmt.Errorf("load: %w", err)
			} else {
				job.Image = img
			}
			out <- job
		}
		close(out)
	}()
	return out
}

func resize(inputChan <-chan Job) <-chan Job {
	out := make(chan Job)
	go func() {
		for job := range inputChan {
			if job.Err == nil {
				job.Image = imageprocessing.Resize(job.Image, DefaultWidth, DefaultHeight)
			}
			out <- job
		}
		close(out)
	}()
	return out
}

func convertToGrayScale(inputChan <-chan Job) <-chan Job {
	out := make(chan Job)
	go func() {
		for job := range inputChan {
			if job.Err == nil {
				job.Image = imageprocessing.Grayscale(job.Image)
			}
			out <- job
		}
		close(out)
	}()
	return out
}

func saveImage(inputChan <-chan Job) <-chan Result {
	out := make(chan Result)
	go func() {
		for job := range inputChan {
			if job.Err != nil {
				out <- Result{Path: job.InputPath, Err: job.Err}
				continue
			}
			err := imageprocessing.WriteImage(job.OutPath, job.Image)
			out <- Result{Path: job.OutPath, Err: err}
		}
		close(out)
	}()
	return out
}
