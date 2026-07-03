FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /go-todo ./cmd/api

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /go-todo /app/go-todo

EXPOSE 8080

CMD ["/app/go-todo"]
