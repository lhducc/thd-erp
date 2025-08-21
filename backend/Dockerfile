## BUILD STAGE##
FROM golang:1.24-alpine AS builder
# FROM 192.168.1.76:5001/product/erp/erp/golang:1.24-alpine AS builder

ENV APPNAME erp-backend

WORKDIR /src

COPY . .

RUN CGO_ENABLED=0 go build -mod=vendor -v -o ./build/$APPNAME ./cmd/api/main.go

## RUN STAGE ##
FROM alpine:3.14
# FROM 192.168.1.76:5001/product/erp/erp/alpine:3.14

ENV APPNAME erp-backend

WORKDIR /app

COPY --from=builder /src/build/$APPNAME .
COPY --from=builder /src/config ./config

# Create SSL directory and copy certificates
RUN mkdir -p /app/ssl
COPY ssl/thd.erp.com.crt /app/ssl/server.crt
COPY ssl/thd.erp.com.key /app/ssl/server.key

# Set proper permissions for SSL files
RUN chmod 644 /app/ssl/server.crt && \
    chmod 600 /app/ssl/server.key

EXPOSE 8080 8443

CMD ["sh", "-c", "./$APPNAME"]
