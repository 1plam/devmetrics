FROM golang:1.23-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/devmetrics

FROM alpine:3.19

WORKDIR /app

RUN adduser -D appuser
USER appuser

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]