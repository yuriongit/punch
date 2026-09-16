# Punch

Punch is a customizable, simple, and lightweight HTTP command-line load testing
tool written in Go. You give it a config file describing the test you want to
run, and it sends requests to your target (and its children (multiple paths)),
and reports how it performed.

## Features

1. Runs HTTP load tests from the command line
2. Configurable load tests with a `punch.json` file.
3. Displays request progress and test results:
   - Individual worker reports successful, fatal, or failed
     requests with response times
   - Workers and global request amounts
   - Formatted measured response times
4. Supports GET requests

## Configuration

Punch uses a `punch.json` file to configure load tests. Additionally, a sample
is included to quickly get started with Punch.

See:

- Configuration Sample: [punch.json](../punch.json)
- Configuration Details: [punch.jsonc](../punch.jsonc)

## Infrastructure

| Layer   | Tool              |
| ------- | ----------------- |
| Main    | Go, Cobra         |
| Tooling | golangci-lint, Go |
| CI      | GitHub Actions    |

## Requirements

- Go 1.27.1 or later

## Quick Start

Clone the repo and build the binary:

```bash
git clone https://github.com/yuriongit/punch.git
cd punch
go build
go install
```

Start Punch:

```bash
punch
```

Run a test; a sample configuration file [`punch.json`](../punch.json) is included:

```bash
punch run ./
```

## Docs

- [planned.md](planned.md)

## License

MIT
