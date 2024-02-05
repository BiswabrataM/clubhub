FROM golang:1.21-alpine

# Work directory
WORKDIR /app

# Installing dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copying all the files
COPY . .

# Starting our application
CMD ["go", "run", "server.go"]

# Exposing server port
EXPOSE 8080