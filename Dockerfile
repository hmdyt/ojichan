FROM golang:1.23.3-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main .

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["/app/main"]
