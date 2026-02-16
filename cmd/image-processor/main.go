package main

import (
	"fmt"
	"os"

	"github.com/revanthvermareddy/image-processing/internal/services"
)

func main() {
	imagePaths := os.Args[1:]
	if len(imagePaths) == 0 {
		imagePaths = []string{
			"./data/images/image-1.jpg",
			"./data/images/image-2.jpg",
		}
		fmt.Println("No image paths provided, using defaults")
	}

	service := services.NewService()
	service.Run(imagePaths)
}
