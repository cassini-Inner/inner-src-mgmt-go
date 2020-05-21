# Start from golang base image
FROM golang:alpine as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .
ENV SERVER_PORT ${SERVER_PORT}
EXPOSE ${SERVER_PORT}
CMD ["./main"]