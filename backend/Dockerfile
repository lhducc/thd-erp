## BUILD STAGE##
FROM golang:1.22-alpine AS builder

ENV APPNAME erp-backend

WORKDIR /src

RUN apk --no-cache add git

COPY go.mod go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -v -o ./build/$APPNAME ./cmd/api/main.go

## RUN STAGE ##
FROM alpine:3.14

ENV APPNAME erp-backend

WORKDIR /run

COPY --from=builder /src/build/$APPNAME .

EXPOSE 8080

CMD ["sh", "-c", "./$APPNAME"]