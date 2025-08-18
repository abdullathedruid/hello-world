# Build stage for CSS
FROM node:18-alpine AS css-builder

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci

# Copy Tailwind config and source files
COPY tailwind.config.js ./
COPY src/ ./src/

# Create static directory and build CSS
RUN mkdir -p static && npm run build-css

# Build stage for Go
FROM golang:1.21-alpine AS go-builder

WORKDIR /app

# Copy go module file
COPY go.mod ./

# Download dependencies (if any)
RUN go mod download

# Copy go source
COPY main.go ./

# Build the Go application
RUN go build -o main .

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