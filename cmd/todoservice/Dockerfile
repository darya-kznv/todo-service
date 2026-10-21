FROM golang:1.27.1-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/cmd/todoservice/exe /app/cmd/todoservice

FROM alpine:3.23 AS runner
WORKDIR /app
COPY --from=builder /app/cmd/todoservice/exe /app
COPY --from=builder /app/public /app/public
ENV PROJECT_ROOT=/app
CMD ["/app/exe"]
