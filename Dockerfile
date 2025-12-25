# syntax=docker/dockerfile:1

FROM golang:1.25-alpine AS build
WORKDIR /src
ENV CGO_ENABLED=0
COPY go.mod go.sum ./
RUN go mod download
COPY cmd cmd
COPY internal internal
RUN go build -trimpath -ldflags='-s -w' -o /out/ip-detect ./cmd/ip-detect

FROM alpine:3.20
WORKDIR /app
RUN addgroup -S app && adduser -S app -G app
COPY --from=build /out/ip-detect /usr/local/bin/ip-detect
EXPOSE 8080
ENV IPD_ADDR=:8080
USER app
ENTRYPOINT ["ip-detect"]
