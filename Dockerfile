# ========================
# 1. Build stage
# ========================
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Переменные для Go сборки, чтобы TMPDIR и GOCACHE не забивали /tmp
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOCACHE=/tmp/go-cache
ENV TMPDIR=/tmp/go-tmp

RUN mkdir -p $GOCACHE $TMPDIR

# Кэшируем зависимости
COPY go.mod go.sum ./
RUN go mod tidy

# Копируем исходники
COPY . .

# Собираем бинарник
RUN go build -ldflags='-w -s' -o main .

# ========================
# 2. Final stage
# ========================
FROM gcr.io/distroless/base-debian11

WORKDIR /

# Копируем бинарник из билд-стадии
COPY --from=builder /app/main ./

# Монтируем сертификаты через volume (не копируем внутрь)
# CMD использует бинарник
EXPOSE 50051
CMD ["/main"]
