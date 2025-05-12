FROM golang:1.24.2 AS build

WORKDIR /app

COPY . .

RUN go build -o vigil

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=build /app/vigil .

COPY .env .

ENTRYPOINT ["./vigil"]
