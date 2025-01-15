# Stage 1: Install dependencies
FROM node:18-alpine as dependencies

WORKDIR /app

# Salin package.json dan package-lock.json untuk cache instalasi npm
COPY package*.json ./

# Install dependencies
RUN npm install --legacy-peer-deps

# Stage 2: Build
FROM node:18-alpine as builder

WORKDIR /app

# Salin seluruh file proyek ke dalam container
COPY . ./

# Jalankan build menggunakan Vite
RUN npx vite build --config apps/frontend/vite.config.ts --outDir dist/apps/frontend

# Stage 3: Production
FROM nginx:alpine

# Salin hasil build ke dalam image nginx
COPY --from=builder /app/dist/apps/frontend /usr/share/nginx/html

# Expose port 80 untuk Nginx
EXPOSE 80

# Start Nginx
CMD ["nginx", "-g", "daemon off;"]
