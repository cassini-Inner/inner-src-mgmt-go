# Start from golang base image
FROM golang:alpine as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .
ARG SERVER_PORT
ENV SERVER_PORT=$SERVER_PORT
RUN echo "SERVER_PORT ${SERVER_PORT}"
EXPOSE ${SERVER_PORT}
CMD ["./main"]