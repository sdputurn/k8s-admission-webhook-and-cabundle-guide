# Use the official Golang image to build the Go application
FROM golang:1.23 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o webhook-server main.go

# # Use a minimal base image to reduce the size of the final image
# FROM alpine:latest

# # Set the Current Working Directory inside the container
# WORKDIR /root/

# # Copy the Pre-built binary file from the previous stage
# COPY --from=builder /app/webhook-server .

# Expose the port on which the webhook server will run
EXPOSE 443

# Command to run the executable
CMD ["./webhook-server"]
