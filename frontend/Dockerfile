## Build stage ##
FROM 10.57.70.41:5001/product/erp/erp/node:24.0.0-slim AS builder
WORKDIR /src

COPY package.json .npmrc .
RUN npm install --legacy-peer-deps --verbose

COPY . .
RUN npm run build


## Run stage ##
FROM 10.57.70.41:5001/product/erp/erp/nginx:slime

COPY --from=builder /src/dist /usr/share/nginx/html

# Override default nginx.conf
COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
