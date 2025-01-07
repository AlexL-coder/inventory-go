# Use the official Golang image to build the app
FROM golang:1.23 AS builder

# Set Go environment variables for static compilation
ENV CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

# Install gcc and musl-dev for CGO support
RUN apt-get update && apt-get install -y gcc musl-dev

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the application source code
COPY . .

# Include the `ent` folder explicitly
COPY ent/ ent/

# Build the Go app
RUN go build -ldflags '-linkmode external -extldflags "-static"' -o main ./cmd/app

# Use a minimal image to run the compiled app
FROM alpine:latest
RUN apk --no-cache add ca-certificates

# Set the working directory in the runtime container
WORKDIR /root/

# Copy the compiled binary from the builder image
COPY --from=builder /app/main .

# Copy the SQLite database file if it exists
COPY --from=builder /app/inventory.db ./

# Include the `ent` folder explicitly in the runtime container
COPY --from=builder /app/ent ent

# Expose the port the app runs on
EXPOSE 8080

# Run the application
CMD ["./main"]