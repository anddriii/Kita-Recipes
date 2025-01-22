# Stage 1: Build
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy source code dan dependensi
COPY . . 
COPY ./cmd/.env /root/.env

# Download dependensi dan build aplikasi
RUN go mod tidy
RUN go build -o main ./cmd

# Stage 2: Production
FROM alpine:latest

WORKDIR /root/

# Copy hasil build dari stage sebelumnya
COPY --from=builder /app/main .

# Expose port yang sesuai
EXPOSE 8081

# Pastikan aplikasi bisa berjalan
CMD ["/root/main"]
