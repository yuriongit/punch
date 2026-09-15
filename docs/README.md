# Punch

Punch is a customizable, simple, and lightweight HTTP load testing CLI-tool.
Available as both a CLI tool. It's language agnostic, as it's meant to be 
used for anything or project.

## Features

...

## Infrastructure

| Layer | Tool |
| --- | --- |
| Main | Go, Cobra |
| Data | Redis |
| Containerzation | Docker |
| CI/CD | GitHub Actions |
| Tooling | golangci-lint, built-in Go tools |

## Running locally

### Clone the repo

```bash
git clone https://github.com/yuriongit/punch.git
cd punch
```

### Start Punch

punch command

```bash
# use of punch command
go build
go install
punch
```

---

Build and run

```bash
# use of local executable
go build 
./punch
# or
go run .
```

## Docs

[todo.md](todo.md)

## Status

Work in progress
