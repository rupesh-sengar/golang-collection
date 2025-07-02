FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN apk update \
 && apk upgrade --no-cache \
 && apk add --no-cache ca-certificates \
 && update-ca-certificates

RUN go mod download
RUN CGO_ENABLED=0 go build -o auth ./cmd/server

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=builder /app/auth /auth
COPY --from=builder /app/.env /.env

EXPOSE 8080
ENTRYPOINT ["/auth"]