FROM golang:1.26-alpine@sha256:ce864e7223ac17b1775e6fd0b4c0db580c2eb50e7953a427916379e4b92a1628

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
