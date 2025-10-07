FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY parse_cobertura.go .
RUN go build -o parse_cobertura parse_cobertura.go

FROM alpine:3.18

WORKDIR /github/workspace
COPY --from=builder /app/parse_cobertura /usr/local/bin/parse_cobertura

ENTRYPOINT ["parse_cobertura"]

