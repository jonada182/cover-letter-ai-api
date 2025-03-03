# Use an official Golang runtime as the base image
FROM golang:1.24.0-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the local code to the container
COPY . .

# Build the Golang application
RUN go build -o api ./cmd/api

# Specify the command to run the executable
CMD ["./api"]
