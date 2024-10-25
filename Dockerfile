# Build stage
FROM golang:1.22-alpine AS build

WORKDIR /app

# Copy Go mod and sum files
COPY go.mod go.sum ./

# Download Go module dependencies
RUN go mod download

# Copy the entire project to the container
COPY . .

# Change the working directory to 'cmd' where the main.go file is located
WORKDIR /app/cmd

# Build the Go binary
RUN go build -o /app/xtrace .

# Final stage - minimal image
FROM alpine:3.16

WORKDIR /app

# Copy the Go binary from the build stage
COPY --from=build /app/xtrace .

# Set the binary to be executable
RUN chmod +x ./xtrace

# Command to run the Go app
CMD ["./xtrace"]

# After building the Go binary, add this line
RUN ls -l /app/xtrace
