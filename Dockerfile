# --- Build stage ---
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Install git (needed for go mod download of some packages)
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./cmd/server

# --- golang-migrate binary stage ---
FROM golang:1.26-alpine AS migrate-builder

RUN apk add --no-cache git
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# --- Final stage ---
FROM alpine:3.20

WORKDIR /app

# ca-certificates needed for any outbound HTTPS calls; netcat for the db-wait loop
RUN apk add --no-cache ca-certificates netcat-openbsd

COPY --from=builder /app/server .
COPY --from=migrate-builder /go/bin/migrate /usr/local/bin/migrate
COPY db/migrations ./db/migrations
COPY entrypoint.sh .

RUN chmod +x entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["./entrypoint.sh"]
