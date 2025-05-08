## Build stage ##
FROM node:24.0.0-alpine3.20 as builder
WORKDIR /src

COPY . .

RUN npm install
RUN npm run build

## Run stage ##
FROM nginx:alpine

COPY --from=builder /src/dist /usr/share/nginx/html

# Override default nginx.conf
COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
