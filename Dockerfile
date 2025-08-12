
FROM golang:1.24-alpine

# Install necessary development tools
RUN apk add --no-cache gcc musl-dev git

# Install templ
RUN go install github.com/a-h/templ/cmd/templ@latest

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Install air for live reloading in development
RUN go install github.com/air-verse/air@latest

# Copy air configuration
COPY .air.toml ./

# Expose the application port
EXPOSE 8080

# Command for development with hot reload
CMD ["air", "-c", ".air.toml"]
