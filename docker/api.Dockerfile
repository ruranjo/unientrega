# -----------------------------
# BUILD STAGE
# -----------------------------
FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

# Copiamos dependencias primero
COPY go.mod go.sum ./
RUN go mod download

# Copiamos TODO el proyecto (desde raíz)
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o api ./cmd/server

# -----------------------------
# RUNTIME STAGE
# -----------------------------
FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/api .

EXPOSE 8080

CMD ["./api"]
