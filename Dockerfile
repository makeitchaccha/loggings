# two stage build
FROM golang:1.23-bookworm AS builder

RUN apt update && apt install -y make && rm -rf /var/lib/apt/lists/*
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make build

# Runner
FROM debian:bookworm-slim
RUN apt update && apt install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /build/build/bot /app/bot
COPY --from=builder /build/build/deploy /app/deploy

CMD [ "/app/bot" ]
