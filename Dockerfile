FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o app .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/app .

COPY --from=builder /app/historico_precos.json .
COPY --from=builder /app/.env .

CMD ["./app"]
