FROM golang:1.23.4 as builder


WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bot

FROM alpine:latest as runner

COPY --from=builder /go/bot /usr/local/bin/bot

ENV SERVER_LISTEN_PORT 5000

ENTRYPOINT ["bot", "serve"]