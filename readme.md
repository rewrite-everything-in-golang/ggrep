# ggrep

`ggrep` is a powerful, parallel grep implementation written in Go. It offers a fast and modular way to search for patterns in files and directories, leveraging Go's concurrency for high performance.

## Features

- **Parallel Search**: Rapidly searches through multiple files using Go's goroutines.
- **Modular Architecture**: Portions of the logic are separated into a reusable `ggrep` package.
- **Full Grep Support**: Supports standard flags including ignore-case, invert-match, line-numbers, recursive search, and more.
- **Colorized Output**: High-visibility output for terminal users.

## Installation

### From Source

Ensure you have [Go](https://go.dev/dl/) installed (version 1.23 or later).

```bash
git clone https://github.com/rewrite-everything-in-golang/ggrep.git
cd ggrep
go build -o ggrep main.go
```

## Usage

```bash
ggrep [OPTIONS] PATTERN [FILE...]
```

### Common Options

- `-i, --ignore-case`: Ignore case distinctions.
- `-v, --invert-match`: Select non-matching lines.
- `-n, --line-number`: Print line numbers of matches.
- `-r, --recursive`: Search directories recursively.
- `-c, --count`: Print only a count of matching lines.
- `-l, --files-with-matches`: Print only names of files with matches.
- `-b, --byte-offset`: Print byte offset of matches.
- `--color`: Enable colorized output.

### Examples

Search for "func" in all Go files recursively with line numbers:
```bash
ggrep -rn "func" .
```

Case-insensitive search for a pattern in a specific file:
```bash
ggrep -i "pattern" main.go
```

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
