FROM golang:1.26-alpine@sha256:ecb67f987aa1a57794c91d19ade7ae0209f2f0f5496aaef9e6ae85ddf92d02ab

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

# Install air for hot reload
RUN go install github.com/air-verse/air@latest

# Copy go mod files
COPY go.mod* go.sum* ./

# Download dependencies
RUN if [ -f go.mod ]; then go mod download; fi

# Copy source code
COPY . .

# Expose port
EXPOSE 8080

# Default command (can be overridden by docker-compose)
CMD ["air", "-c", ".air.toml"]
