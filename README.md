# Punch

Punch is a HTTP load tester available, as both a CLI tool and web app.
Built to deepen my understanding of Docker, GitHub Actions, Go,
and concurrent systems design.

![Punch UI Preview](.github/images/preview.png)

## Features

...

## Technology

**Web App**

| Layer | Tools |
| --- | --- |
| **Frontend** | React, TypeScript, Vite, Bun |
| **API** | Go, Gin, Docker, Redis |
| **CI/CD & Hosting** | GitHub Actions, Vercel, Railway & GCP |

**CLI**

| Layer | Tools |
| --- | --- |
| **API** | Go, Docker, Redis |

## Running locally

### Clone the repo

```bash
git clone https://github.com/yuriongit/punch.git
cd punch
```

### Start Frontend

```bash
cd web
bun install --frozen-lockfile
bun run build 
bun run preview
```

### Start API

```bash
go run .
```

## Docs

...

## Status

Work in progress
