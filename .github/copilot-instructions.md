# Copilot Instructions for image-processing

## Project Overview

This is a Go image-processing pipeline that uses goroutines and channels to
process images concurrently. Each stage (load → resize → grayscale → save) runs
in its own goroutine and communicates via typed channels.

## Architecture

```
cmd/image-processor/main.go   – CLI entry point
internal/services/             – Pipeline orchestration (unexported stages)
pkg/imageprocessing/           – Reusable image I/O and transformation functions
data/images/                   – Sample input images
data/images/output/            – Processed output images (git-ignored)
```

### Key Design Decisions
- **Channel pipeline pattern**: each stage is a function that returns a
  receive-only channel, enabling composable concurrent pipelines.
- **Error propagation**: errors are carried in the `Job.Err` field so
  downstream stages can skip broken jobs without panicking.
- **`pkg/` vs `internal/`**: image functions in `pkg/` are importable by
  external consumers; pipeline wiring in `internal/` is implementation detail.

## Conventions
- All exported functions and types must have GoDoc comments.
- Functions in `pkg/imageprocessing` return `error` instead of panicking.
- Use `t.TempDir()` in tests for temporary files.
- Run `make lint` before committing.

## Common Tasks
```bash
make build       # compile binary to ./bin/
make run         # build and run with default images
make test        # run all tests with race detector
make lint        # run go vet
make clean       # remove build artifacts
```
