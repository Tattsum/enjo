FROM golang:1.26-alpine@sha256:3082400e369fa24d5fc60bca20edab3f6d604e0c5a690ec66b295eff4dd87ade

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
