FROM golang:1.26.2-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /app/todoky ./cmd/todoky

FROM alpine:3.23

WORKDIR /app

COPY --from=builder /app/todoky /app/todoky

CMD ["/app/todoky"]