# syntax=docker/dockerfile:1
FROM quay.io/projectquay/golang:1.22

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags "-X main.version=$(git log -1 --format=%cd --date=format:'%Y%m%d')" -o /app/cmd/wildlifenl cmd/main.go

EXPOSE 8080

CMD ["/app/cmd/wildlifenl"]