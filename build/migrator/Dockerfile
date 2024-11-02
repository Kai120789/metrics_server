FROM golang:1.23.1-alpine

WORKDIR /app

COPY ../go.mod ../go.sum ./

RUN go mod download

COPY . .

RUN go build -o ./goapp2 cmd/migrator/migrator.go

CMD ["./goapp2"]