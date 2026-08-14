# uncomment-cli

A fast and flexible CLI tool that removes **all** comments from Go source files — one file at a time or many files in parallel.

[![Go version](https://img.shields.io/github/go-mod/go-version/AlGrushino/uncomment_cli)](https://go.dev/)
[![Go Reference](https://pkg.go.dev/badge/github.com/AlGrushino/uncomment_cli.svg)](https://pkg.go.dev/github.com/AlGrushino/uncomment_cli)
[![CI Pipeline](https://github.com/AlGrushino/uncomment_cli/actions/workflows/pipeline.yaml/badge.svg)](https://github.com/AlGrushino/uncomment_cli/actions)
[![License](https://img.shields.io/github/license/AlGrushino/uncomment_cli)](LICENSE)

## Features

- **Strips every kind of comment** — line (`//`), block (`/* */`) and doc comments — in a single AST pass.
- **Safe on strings** — comment-like text inside string literals (e.g. `fmt.Println("// not a comment")`) is left untouched.
- **Batch mode** — `uncomment-many` processes any number of files concurrently using a worker pool sized to the number of CPU cores.
- **Path safety guard** — only files inside your home or temp directory are allowed, everything else is rejected.
- **Clean output** — the result is passed through `go/format`, so the code stays formatted and valid.
- **Go standard library only** — parsing is built on `go/parser`, `go/ast` and `go/format`.

## Requirements

- Go 1.25.2 or newer
- GNU Make (optional, but recommended)

## Installation

Clone the repository and build the binary:

```bash
git clone https://github.com/AlGrushino/uncomment_cli.git
cd uncomment_cli
make build
```

Or install it directly into your `$GOBIN`:

```bash
go install github.com/AlGrushino/uncomment_cli@latest
```

The binary will be placed at `bin/uncomment-cli`.

You can also run it without building:

```bash
go run . <command> [flags]
```

## Usage

Remove comments from a single file:

```bash
uncomment-cli uncomment -p ./main.go
# or, using the long form
uncomment-cli uncomment --path ./main.go
```

Remove comments from several files at once (the `-p` flag can be repeated):

```bash
uncomment-cli uncomment-many -p ./main.go -p ./internal/app.go -p ./pkg/util.go
```

### Example

Before:

```go
// main package comment
package main

import "fmt"

// prints a greeting
func greet() {
    fmt.Println("Hello, world! // this is a string, not a comment")
    /* an old
       block comment */
}
```

After running `uncomment-cli uncomment -p main.go`:

```go
package main

import "fmt"

func greet() {
	fmt.Println("Hello, world! // this is a string, not a comment")
}
```

## Command reference

| Command | Description | Flags |
| --- | --- | --- |
| `uncomment` | Deletes all comments from a single `*.go` file | `-p, --path string` — path to the file |
| `uncomment-many` | Deletes all comments from any number of `*.go` files (concurrently) | `-p, --pathes stringArray` — paths to the files, may be repeated |

Both commands print a message per file and report errors without stopping the batch.

## How it works

1. The file is parsed into an AST with `go/parser` (with comment nodes enabled).
2. All comment groups are dropped, and every node's `Doc`/`Comment` fields are cleared.
3. The cleaned AST is formatted back to source code with `go/format`.
4. The resulting text is written back to the original file — but only after the path passes the home/temp-directory safety check in `pkg/path`.

## Project layout

```
.
├── cmd/                   # Cobra commands (root, uncomment, uncomment-many)
├── internal/
│   ├── fs/                # Filesystem abstraction
│   ├── service/           # Core uncommenting logic (single & many)
│   └── utils/             # Filename / path helpers
├── pkg/path/              # Allowed-path validation
├── main.go                # Entry point
├── Makefile               # Build, test, lint, coverage tasks
└── .github/workflows/     # CI pipeline
```

## Development

| Command | Description |
| --- | --- |
| `make build` | Build the CLI binary |
| `make run` | Run the CLI with `go run` |
| `make test` | Run all tests |
| `make test.coverage` | Run tests and print a coverage report |
| `make test.html` | Generate an HTML coverage report (`coverage.html`) |
| `make fmt` | Format the code with `go fmt` |
| `make lint` | Run `go vet`, `golangci-lint` and `gosec` |
| `make deps` | Tidy module dependencies |
| `make clean` | Remove build artifacts and coverage files |
| `make help` | List all available Makefile targets |

## CI

Every push or pull request to `main` / `develop` runs the pipeline in `.github/workflows/pipeline.yaml`:

1. `fmt` — checks that the code is formatted with `go fmt`.
2. `lint` — runs `golangci-lint`.
3. `test` — runs `make test`.

## License

[MIT](LICENSE) © 2026 AlGrushino
