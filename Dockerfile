# Use the official Go image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod .
COPY go.sum .

# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o main .

# Expose the port your application will run on
EXPOSE 8080

# Command to run your application
CMD ["./main"]
