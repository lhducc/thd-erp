## Build stage ##
# FROM 192.168.1.76:5001/product/erp/erp/node:24.0.0-slim AS builder
FROM node:24.0.0-slim AS builder
WORKDIR /src

# Accept build arguments
ARG VITE_API_URL
ENV VITE_API_URL=$VITE_API_URL

# Copy package files first for better Docker layer caching
COPY package*.json ./
RUN npm install --legacy-peer-deps

# Copy source code and build
COPY . .
RUN npm run build


## Run stage ##
# FROM 192.168.1.76:5001/product/erp/erp/nginx:1.27
FROM nginx:1.27

COPY --from=builder /src/dist /usr/share/nginx/html

# Create SSL directory and copy certificates
RUN mkdir -p /etc/nginx/ssl
COPY ssl/thd.erp.com.crt /etc/nginx/ssl/server.crt
COPY ssl/thd.erp.com.key /etc/nginx/ssl/server.key

# Set proper permissions for SSL files
RUN chmod 644 /etc/nginx/ssl/server.crt && \
    chmod 600 /etc/nginx/ssl/server.key

# Override default nginx.conf
COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 80 443

CMD ["nginx", "-g", "daemon off;"]
