# syntax=docker/dockerfile:1

FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY api/go.mod api/go.sum ./
RUN go mod download
COPY api .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main .

FROM alpine:3.22.0 AS runner
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations
EXPOSE 3000

CMD ["./main"]
