## Build stage ##
# FROM node:24.0.0-slim AS builder
FROM 192.168.1.76:5001/product/erp/erp/node:24.0.0-slim AS builder
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