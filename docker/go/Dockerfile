FROM golang:1.26-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /url-shortener ./cmd/url-shortener

FROM scratch
WORKDIR /app
COPY --from=builder /url-shortener /app/
COPY config/local.yaml /app/config/config.yaml

EXPOSE 8080
ENTRYPOINT ["/app/url-shortener"]
