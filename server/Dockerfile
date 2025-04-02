FROM golang:1.24.1-alpine3.20 AS builder

WORKDIR /app

COPY . .

RUN go build -ldflags '-w -s -extldflags "-static"' -a  .

FROM alpine:3.21.3

WORKDIR /app

COPY --from=builder /app ./

EXPOSE 8000

CMD ["./server"]
