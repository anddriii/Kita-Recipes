# Gunakan image node
FROM node:18

# Set working directory
WORKDIR /app

# Copy package.json dan package-lock.json dari direktori frontend
COPY ../../package*.json ./

# Install dependencies
RUN npm install --legacy-peer-deps

# Copy seluruh kode frontend
COPY ../../apps/frontend ./

# Build aplikasi frontend
RUN npm run build:frontend

# Gunakan image nginx untuk serving file statis
FROM nginx:alpine
COPY --from=builder /app/dist/apps/frontend /usr/share/nginx/html

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
