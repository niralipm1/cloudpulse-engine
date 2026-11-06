# Build Stage
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download || true

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /cloudpulse-engine main.go

# Production Stage
FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

COPY --from=builder /cloudpulse-engine .

EXPOSE 9092
ENTRYPOINT ["./cloudpulse-engine"]
