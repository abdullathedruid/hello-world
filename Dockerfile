# Build stage for CSS
FROM node:18-alpine AS css-builder

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies with cache mount
RUN --mount=type=cache,target=/root/.npm \
    npm ci

# Copy Tailwind config and source files
COPY tailwind.config.js ./
COPY src/ ./src/

# Create static directory and build CSS
RUN mkdir -p static && npm run build-css

# Build stage for Go
FROM golang:1.24-alpine AS go-builder

WORKDIR /app

# Copy go module files
COPY go.mod go.sum ./

# Download dependencies with cache mount
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy go source
COPY main.go ./

# Build the Go application with cache mount
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o main .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Go binary
COPY --from=go-builder /app/main .

# Copy static files from CSS build stage
COPY --from=css-builder /app/static ./static/

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]