FROM golang:1.24.3-alpine3.21 AS builder

WORKDIR /app

COPY . .

RUN go build -ldflags '-w -s -extldflags "-static"' -a -o gitstats .

FROM alpine:3.21.3

WORKDIR /app

COPY --from=builder /app ./

EXPOSE 8000

CMD ["./gitstats"]
