FROM golang:1.24-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/waitlistfox .

FROM alpine:3.22

WORKDIR /app

RUN addgroup -S app && adduser -S -G app app \
    && apk add --no-cache ca-certificates tzdata \
    && mkdir -p /app/config /app/internal/i18n /app/log

COPY --from=builder /out/waitlistfox /app/waitlistfox
COPY internal/i18n /app/internal/i18n
COPY static /app/static
COPY config/config.example.json /app/config/config.example.json

RUN chown -R app:app /app

USER app

ENV CONFIG_PATH=/app/config/config.json

EXPOSE 8081

CMD ["./waitlistfox"]
