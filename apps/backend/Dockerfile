# Gunakan image golang
FROM golang:1.20

# Set working directory
WORKDIR /app

# Copy file go.mod dan go.sum
COPY ./go.mod ./go.sum ./

# Download dependencies
RUN go mod download

# Copy seluruh kode backend
COPY ./ ./

# Build binary file
RUN go build -o main ./cmd

# Expose port untuk backend
EXPOSE 8080

# Command untuk menjalankan aplikasi
CMD ["./main"]
