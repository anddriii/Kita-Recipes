# Stage 1: Build Backend (Golang)
FROM golang:1.21 as backend-builder

WORKDIR /app/backend

# Salin go.mod dan go.sum
COPY apps/backend/go.mod apps/backend/go.sum ./

# Unduh dependencies untuk backend
RUN go mod tidy

# Salin kode backend
COPY apps/backend .

# Set working directory ke folder 'cmd'
WORKDIR /app/backend/cmd

# Build aplikasi backend
RUN go build -o /app/backend/main .

# Stage 2: Final Image
FROM node:16

WORKDIR /app

# Salin file package.json dan package-lock.json
COPY package.json package-lock.json ./

# Install dependencies root (untuk menjalankan script npm)
RUN npm install

# Salin hasil build backend
COPY --from=backend-builder /app/backend/main ./backend/

# Salin semua kode untuk runtime
COPY . .

# Reset cache NX sebelum menjalankan
RUN npx nx reset

# Jalankan frontend dan backend secara bersamaan
CMD ["npm", "run", "start"]

# Expose ports untuk backend dan frontend
EXPOSE 8080 4200
