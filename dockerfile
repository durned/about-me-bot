# Use the official Golang image to build the application
FROM golang:latest AS builder

# Create a workdir inside the container
WORKDIR /app/about-me-bot
# Switch
WORKDIR /app

COPY ./configs ./configs

# Copy the Go Modules manifests
COPY ./about-me-bot/go.mod ./about-me-bot/go.sum ./

# Download dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY ./about-me-bot ./about-me-bot

WORKDIR /app/about-me-bot/cmd

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o /my-go-app

# Command to run the executable
CMD ["/my-go-app"]
