package services

import (
	"fmt"
	"image"
	"strings"

	"github.com/revanthvermareddy/image-processing/pkg/imageprocessing"
)

type Service interface {
	Run(imagePaths []string)
}

type ServiceImpl struct{}

func NewService() Service {
	return &ServiceImpl{}
}

func (s *ServiceImpl) Run(imagePaths []string) {

	channel1 := loadImage(imagePaths)
	channel2 := resize(channel1)
	channel3 := convertToGrayScale(channel2)
	writeResults := saveImage(channel3)

	for success := range writeResults {
		if success {
			fmt.Println("Image saved successfully")
		} else {
			fmt.Println("Image not saved")
		}
	}
}

type Job struct {
	InputPath string
	Image     image.Image
	OutPath   string
}

func loadImage(paths []string) <-chan Job {
	// create a new channel
	out := make(chan Job)

	// create a goroutine to process the images
	go func() {
		// For each image path, create a job and add it to the out channel
		for _, path := range paths {
			job := Job{
				InputPath: path,
				OutPath:   strings.Replace(path, "images/", "images/output/", 1),
			}
			job.Image = imageprocessing.ReadImage(path)
			out <- job
		}
		// Close the channel when done
		close(out)
	}()

	return out
}

func resize(inputChan <-chan Job) <-chan Job {
	// create a new channel
	out := make(chan Job)

	// create a goroutine to process the images
	go func() {
		// For each input job, update the job by updating the image after resizing and add it to the out channel
		for job := range inputChan { // Read from the input channel
			job.Image = imageprocessing.Resize(job.Image)
			out <- job
		}
		// Close the channel when done
		close(out)
	}()

	return out
}

func convertToGrayScale(inputChan <-chan Job) <-chan Job {
	// create a new channel
	out := make(chan Job)

	// create a goroutine to process the images
	go func() {
		// For each input job, update the job by converting the image to grayscale and add it to the out channel
		for job := range inputChan { // Read from the input channel
			job.Image = imageprocessing.Grayscale(job.Image)
			out <- job
		}
		// Close the channel when done
		close(out)
	}()

	return out
}

func saveImage(inputChan <-chan Job) <-chan bool {
	// create a new channel
	out := make(chan bool)

	// create a goroutine to process the images
	go func() {
		// For each input job, save the image to file system and add the status to the out channel
		for job := range inputChan { // Read from the input channel
			imageprocessing.WriteImage(job.OutPath, job.Image)
			out <- true
		}
		// Close the channel when done
		close(out)
	}()

	return out
}
