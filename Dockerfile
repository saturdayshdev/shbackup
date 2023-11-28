FROM golang:1.21.4-alpine3.18 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o bin ./cmd

FROM alpine:3.14.2

WORKDIR /app

COPY --from=builder /app/bin .

CMD ["./bin"]