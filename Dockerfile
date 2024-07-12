# syntax=docker/dockerfile:1
FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY cmd/*.go cmd/

RUN go build -o /app/wildlifenl cmd/*.go

EXPOSE 8080

CMD ["/app/wildlifenl"]