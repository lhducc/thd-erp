## BUILD STAGE##
FROM 10.57.70.41:5001/product/erp/erp/golang:1.22-alpine AS builder

ENV APPNAME erp-backend

WORKDIR /src

RUN apk --no-cache add git

COPY go.mod go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -v -o ./build/$APPNAME ./cmd/api/main.go

## RUN STAGE ##
FROM 10.57.70.41:5001/product/erp/erp/alpine:3.14

ENV APPNAME erp-backend

WORKDIR /app

COPY --from=builder /src/build/$APPNAME .

EXPOSE 8080

CMD ["sh", "-c", "./$APPNAME"]
