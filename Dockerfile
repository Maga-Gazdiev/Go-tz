FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY . .

RUN CGO_ENABLED=0 go build -o /server ./cmd/api

FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /server .
COPY docs ./docs
EXPOSE 8099

CMD ["./server"]
