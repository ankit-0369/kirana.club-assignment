# Use an official Go runtime as a base image
FROM golang:1.20-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application files
COPY . .

# Build the application
RUN go build -o app .

# Use a lightweight image to run the application
FROM alpine:latest

# Set the working directory
WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/app .

# Expose the port that the app runs on
EXPOSE 8080

# Run the app binary
CMD ["./app"]
