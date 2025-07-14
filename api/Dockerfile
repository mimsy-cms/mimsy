# syntax=docker/dockerfile:1

FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main .

FROM alpine:3.22.0 AS runner
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 3000

CMD ["./main"]
