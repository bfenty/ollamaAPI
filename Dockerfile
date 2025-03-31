# Use the official Golang image as the base image
FROM golang:1.20-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to the container
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code to the container
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

# Use a smaller base image for the final container
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /root/

# Copy the built application from the builder stage
COPY --from=builder /app/app .

# Expose the port the application is running on (if applicable)
EXPOSE 8080

# Run the application
CMD ["./app"]
