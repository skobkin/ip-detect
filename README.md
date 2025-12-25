# Minimalistic IP detector

[![Build Status](https://ci.skobk.in/api/badges/skobkin/ip-detect/status.svg)](https://ci.skobk.in/skobkin/ip-detect)

A tiny Go service that returns your remote address metadata in HTML, JSON, or plain text. It mirrors the old PHP utility while adding reverse DNS lookups, request logging, and a refreshed UI.

## Requirements
- Go 1.25 or newer

## Run locally
```bash
# Serve on the default :8080
go run ./cmd/ip-detect

# Override settings with environment variables
IPD_ADDR=":9090" \
IPD_TRUSTED_SUBNETS="10.0.0.0/8,192.168.0.0/16" \
    go run ./cmd/ip-detect
```

## Configuration
All knobs are exposed via environment variables prefixed with `IPD_`. Common options:

| Variable                         | Default | Purpose                                                                                                           |
|----------------------------------|---------|-------------------------------------------------------------------------------------------------------------------|
| `ADDR`                           | `:8080` | Bind address/port.                                                                                                |
| `READ_TIMEOUT` / `WRITE_TIMEOUT` | `5s`    | HTTP read/write limits.                                                                                           |
| `SHUTDOWN_TIMEOUT`               | `10s`   | Graceful shutdown timeout.                                                                                        |
| `TRUST_FORWARDED`                | `true`  | Whether to honor `X-Forwarded-For` / `X-Real-IP`.                                                                 |
| `TRUSTED_SUBNETS`                | ``      | Comma-separated CIDRs required to trust proxy headers (empty = trust every proxy when `TRUST_FORWARDED` is true). |
| `RESOLVE_PTR`                    | `true`  | Resolve PTR records for the detected IP.                                                                          |
| `RESOLVE_TIMEOUT`                | `500ms` | Reverse DNS lookup timeout per request.                                                                           |
| `INCLUDE_UA`                     | `true`  | Attach the `User-Agent` header to responses.                                                                      |
| `INCLUDE_TS`                     | `true`  | Emit the current UTC timestamp.                                                                                   |
| `LOG_LEVEL`                      | `info`  | One of `debug`, `info`, `warn`, `error`.                                                                          |
| `LOG_FORMAT`                     | `text`  | `text` or `json` output.                                                                                          |

## Docker
An Alpine-based multi-stage build is provided:
```bash
docker build -t skobkin/ip-detect .
docker run --rm -p 8080:8080 skobkin/ip-detect
```
