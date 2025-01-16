# Stage 1: Install dependencies
FROM node:18 as dependencies


WORKDIR /app

# Salin package.json dan package-lock.json untuk cache instalasi npm
COPY package*.json ./

# Hapus node_modules dan package-lock.json sebelum instalasi ulang
RUN rm -rf node_modules package-lock.json && yarn install


# Install dependencies
RUN yarn install @rollup/rollup-linux-x64-gnu --legacy-peer-deps



# Stage 2: Build
FROM node:18 as builder


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
