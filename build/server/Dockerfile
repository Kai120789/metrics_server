FROM golang:1.23.1-alpine

WORKDIR /app

COPY ../go.mod ../go.sum ./

RUN go mod download

COPY . .

RUN go build -o ./goapp cmd/server/server.go

CMD ["./goapp"]