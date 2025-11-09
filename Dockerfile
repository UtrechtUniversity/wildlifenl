# syntax=docker/dockerfile:1
FROM quay.io/projectquay/golang:1.22

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY database /database

COPY templates /templates

RUN go build -ldflags "-X main.version=$(git log -1 --format=%cd --date=format:'%Y%m%d')" -o /app/wildlifenl cmd/main.go

EXPOSE 8080

CMD ["/app/wildlifenl"]