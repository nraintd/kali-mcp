FROM golang:1.26 AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN VERSION="$(awk -F= '/^VERSION=/{print $2}' VERSION.env | tr -d '\r')" && \
  go tool task build TARGET_OS=linux TARGET_ARCH=amd64 VERSION=${VERSION}

FROM kalilinux/kali-rolling:latest

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y --no-install-recommends \
  ca-certificates \
  bash \
  kali-linux-default \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /src/bin/kali-mcp-*-linux-amd64 /usr/local/bin/kali-mcp

EXPOSE 7075 7076

ENTRYPOINT ["/usr/local/bin/kali-mcp"]
