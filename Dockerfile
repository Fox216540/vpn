# ========================
# 1. Build stage
# ========================
FROM golang:1.21 AS builder

WORKDIR /build

# Кэширование зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем и собираем
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags='-w -s' -o main .

# ========================
# 2. Final stage
# ========================
FROM gcr.io/distroless/base-debian11

WORKDIR /
COPY --from=builder /build/main .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 50051

CMD ["/main"]