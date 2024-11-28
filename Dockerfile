# Base image for Go
FROM golang:alpine

LABEL maintainer="online-questionnaire"

# Install necessary tools including Delve
RUN apk add --no-cache git bash gcc musl-dev && \
    go install github.com/go-delve/delve/cmd/dlv@latest

RUN mkdir /app
WORKDIR /app

# Copy application files
COPY . .
COPY .env .

# Setup Go environment
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download
# Expose ports for the app and delve
EXPOSE 8080 40000

# Use Delve to run the application in debug mode
CMD ["go","run","./cmd/app/main.go"]