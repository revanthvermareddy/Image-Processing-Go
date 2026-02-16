# Image Processing Pipeline

A concurrent image-processing pipeline written in Go. It uses goroutines and channels to process multiple images in parallel through a series of stages.

## Pipeline Stages

```
Load Image → Resize (500×500) → Convert to Grayscale → Save to Disk
```

Each stage runs in its own goroutine and communicates with the next via a typed channel, following Go's [pipeline concurrency pattern](https://go.dev/blog/pipelines).

## Project Structure

```
.
├── cmd/image-processor/       # CLI entry point
│   └── main.go
├── internal/services/         # Pipeline orchestration (Job, Result types)
│   ├── services.go
│   └── services_test.go
├── pkg/imageprocessing/       # Reusable image I/O and transformations
│   ├── image_processing.go
│   └── image_processing_test.go
├── data/images/               # Sample input images
│   └── output/                # Processed output (git-ignored)
├── Makefile
├── go.mod
└── go.sum
```

| Directory | Purpose |
|---|---|
| `cmd/` | Application entry points |
| `internal/` | Private application code — not importable by other modules |
| `pkg/` | Public library code — safe to import in other projects |
| `data/` | Sample data and output |

## Getting Started

### Prerequisites

- [Go 1.25+](https://go.dev/dl/)

### Build & Run

```bash
# Build the binary
make build

# Run with default sample images
make run

# Or pass custom image paths
./bin/image-processor path/to/image1.jpg path/to/image2.jpg
```

### Testing

```bash
make test
```

### Linting

```bash
make lint
```

## How It Works

1. **`pkg/imageprocessing`** provides four pure functions:
   - `ReadImage(path) → (image.Image, error)` — decodes a JPEG file
   - `WriteImage(path, img) → error` — encodes and writes a JPEG (creates parent dirs)
   - `Resize(img, w, h) → image.Image` — scales using Lanczos3 interpolation
   - `Grayscale(img) → image.Image` — converts to 8-bit grayscale

2. **`internal/services`** wires these functions into a channel pipeline:
   - `loadImage` — reads files and sends `Job` structs into a channel
   - `resize` — reads jobs, resizes images, forwards them
   - `convertToGrayScale` — reads jobs, converts to grayscale, forwards them
   - `saveImage` — reads jobs, writes files, sends `Result` structs

3. **`cmd/image-processor`** accepts image paths as CLI arguments (or uses defaults) and calls `services.Run()`.

Errors are propagated through the pipeline via the `Job.Err` field — no panics.

## License

This project is for educational purposes.
