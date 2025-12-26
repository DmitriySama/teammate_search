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
COPY --from=builder src/internal/frontend /app/internal/frontend
ENV PORT=3000
ENV CONFIG_PATH=/app/config.yml
ENV FRONTEND_PATH=/app/internal/frontend
EXPOSE 3000

CMD ["/app/teammate-search"]
