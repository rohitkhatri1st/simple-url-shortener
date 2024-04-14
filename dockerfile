# Use an official Golang image for the appropriate platform
FROM golang:latest

# Set the working directory within the container
WORKDIR /app

# Copy the application source code into the container
COPY . .

# Perform "go mod tidy" to manage dependencies
RUN go mod tidy

# Build Application
RUN go build main.go

# Expose port
EXPOSE 8001

# Build and run your Go application
ENTRYPOINT ["./main"]
