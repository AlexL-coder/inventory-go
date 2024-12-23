# Use the official Golang image to build the app
FROM golang:1.23 as builder

# Set Go environment variables
ENV CGO_ENABLED=0 \
    GOOS=linux

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the application source code
COPY . .

# Build the Go app
RUN go build -o main ./cmd/app

# Use a minimal image to run the compiled app
FROM alpine:latest
RUN apk --no-cache add ca-certificates

# Set the working directory in the runtime container
WORKDIR /root/

# Copy the compiled binary from the builder image
COPY --from=builder /app/main .

# Copy the SQLite database file if it exists
COPY --from=builder /app/inventory.db .

# Expose the port the app runs on
EXPOSE 8080

# Run the application
CMD ["./main"]
