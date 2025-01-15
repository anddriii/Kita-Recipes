# Stage 1: Install dependencies
FROM golang:1.19-alpine as dependencies

WORKDIR /app

# Copy file go.mod dan go.sum
COPY ./go.mod ./go.sum ./

# Download dependensi
RUN go mod download

# Stage 2: Build
FROM golang:1.19-alpine as builder

WORKDIR /app

# Copy semua file dari stage sebelumnya
COPY --from=dependencies /go/pkg /go/pkg

# Copy sisa source code
COPY . .

# Build aplikasi
RUN go build -o main ./cmd

# Stage 3: Production
FROM alpine:latest

WORKDIR /root/

# Copy hasil build dari stage sebelumnya
COPY --from=builder /app/main .

# Expose port 8080 untuk aplikasi Go
EXPOSE 8080

CMD ["./main"]
