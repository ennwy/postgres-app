FROM golang:alpine

LABEL maintainer="Agus Wibawantara"

RUN apk update && apk add --no-cache git && apk add --no-cach bash && apk add build-base

# Setup folders
RUN mkdir /app
WORKDIR /app

ENV GOPATH=/

COPY ./ ./
COPY .env ./
COPY go.mod go.sum ./

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...
RUN go mod download
RUN go build -o main main.go

EXPOSE 8080

CMD ["./main"]