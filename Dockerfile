FROM golang:1.24-alpine

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o main cmd/web/main.go

# Copy templates
COPY templates ./templates

# Expose the port
EXPOSE 8080

# Run the application
CMD ["./main"] 