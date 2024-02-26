# Use the official golang image as base image
FROM golang:1.20-alpine AS build

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY src/web/main.go main.go

# Build the Go application
RUN go build -o web

# Use a lightweight Alpine image to package the application
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /root

# Copy the binary from the build stage to the final stage
COPY --from=build /app/web .
COPY src/web/index.html .

# Command to run the application
CMD ["./web"]
