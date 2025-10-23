package main

import (
	"github.com/revanthvermareddy/image-processing/internal/services"
)

func main() {
	imagePaths := []string{
		"./data/images/image-1.jpg",
		"./data/images/image-2.jpg",
	}

	service := services.NewService()
	service.Run(imagePaths)
}
