FROM golang:1.25 AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/teammate-search ./cmd/app

FROM alpine:3.19

RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /bin/teammate-search /app/teammate-search
COPY config.yml /app/config.yml
ENV PORT=8063
ENV CONFIG_PATH=/app/config.yml
EXPOSE 8063

CMD ["/app/teammate-search"]

