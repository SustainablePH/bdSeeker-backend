# Start from official Go image
FROM golang:1.22-alpine

# Set working directory inside container
WORKDIR /app

# Copy Go modules and download them
COPY go.mod ./
RUN go mod download

# Copy rest of the code
COPY . .

# Build the app
RUN go build -o server .

# Expose port and run
EXPOSE 8080
CMD ["./server"]
