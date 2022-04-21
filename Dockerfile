FROM golang:1.17.0-alpine3.14 as builder

WORKDIR /go/src/github.com/joleques/go-redis-stream

COPY . .

RUN go mod download

RUN go build -ldflags "-s -w" consumer/main.go

FROM alpine:3.14

WORKDIR /app

COPY --from=builder /go/src/github.com/joleques/go-redis-stream/main .

CMD [ "./main" ]
