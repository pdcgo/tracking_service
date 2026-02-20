# Use the official Go image
FROM golang:1.24.4 AS builder

WORKDIR /app
COPY go.mod go.sum ./


COPY . .
RUN go mod tidy
RUN go build -o production ./cmd/production

# Small runtime image
FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /app/production .

ENV PORT=8080

ENTRYPOINT ["/app/production"]