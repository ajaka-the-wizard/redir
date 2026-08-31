FROM golang:1.25-alpine AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/redir ./cmd/api

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
RUN addgroup -S redir && adduser -S -G redir redir
WORKDIR /app

COPY --from=builder /bin/redir /app/redir
RUN chown redir:redir /app/redir

EXPOSE 5000

USER redir

ENTRYPOINT ["/app/redir"]
