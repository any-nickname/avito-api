# Шаг 1: Кеширование модулей
FROM golang:1.20 as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Шаг 2: Сборка
FROM golang:1.20 as builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -tags migrate -o /bin/app ./cmd/app

# Шаг 3: Финал
FROM scratch
COPY --from=builder /app/config /config
COPY --from=builder /app/migrations /migrations
COPY --from=builder /bin/app /app
CMD ["/app"]