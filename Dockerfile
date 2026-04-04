ARG GO_VERSION=1.26.1
ARG ALPINE_VERSION=3.22

FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /src

RUN apk add --no-cache ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/api ./cmd/api

FROM alpine:${ALPINE_VERSION} AS runtime

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app \
    && adduser -S -G app app

COPY --from=builder /app/api /app/api

USER app

EXPOSE 8080

CMD ["./api"]
