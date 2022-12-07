FROM golang:1.19

WORKDIR /Users/etaurnv/go/src/github.com/chumvan/confdb

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -o main ./server/main.go

EXPOSE 8080

CMD ["./main"]

