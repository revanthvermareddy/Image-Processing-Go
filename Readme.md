# Image-Processing-Go

This project aims to build an efficient go-pipeline for image processing. It utilizes the go concurrency pattern to achieve high performance.
The pipeline consists of multiple stages, each responsible for a specific image processing task.

Stages:
- Reading the image from a file
- Compressing the image
- Gray scaling the image
- Writing image to a file

### Design Pattern

I'm using the concurrency pattern here which uses channels to communicate between the stages. Each stage is implemented as a goroutine that reads from an input channel, processes the data, and writes to an output channel.
