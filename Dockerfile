# Build stage for frontend assets
FROM node:18-alpine AS frontend-builder

WORKDIR /app

# Copy package files
COPY package*.json ./
COPY tsconfig.json ./

# Install dependencies with cache mount
RUN --mount=type=cache,target=/root/.npm \
    npm ci

# Copy Tailwind config and source files
COPY tailwind.config.js ./
COPY src/ ./src/

# Create static directory and build all frontend assets
RUN mkdir -p static/css static/js && npm run build

# Build stage for Go
FROM golang:1.24-alpine AS go-builder

WORKDIR /app

# Accept build argument for commit hash
ARG COMMIT_HASH
ENV COMMIT_HASH=${COMMIT_HASH:-unknown}

# Copy go module files
COPY go.mod go.sum ./

# Download dependencies with cache mount
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy go source files
COPY main.go ./
COPY config/ ./config/
COPY handlers/ ./handlers/
COPY middleware/ ./middleware/
COPY models/ ./models/
COPY routes/ ./routes/
COPY services/ ./services/

# Build the Go application with cache mount
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -ldflags "-X main.CommitHash=${COMMIT_HASH}" -o main .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Go binary
COPY --from=go-builder /app/main .

# Copy static files from frontend build stage
COPY --from=frontend-builder /app/static ./static/

# Copy template files
COPY templates/ ./templates/

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]