# Development Dockerfile with hot reload support
FROM golang:1.24-alpine

WORKDIR /app

# Install air for hot reload and other tools
RUN go install github.com/air-verse/air@latest

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code (will be overridden by volume mount)
COPY . .

EXPOSE 8080

# Air will watch for changes and rebuild automatically
CMD ["air", "-c", ".air.toml"]
