FROM golang:1.24.3 AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# کپی go.mod از root directory
COPY ../../../go.mod ../../../go.sum ./

RUN go mod download

# کپی shared modules
COPY ../../shared ./shared

# کپی user service code
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /root/

COPY --from=builder /app/main .

RUN chown appuser:appgroup main

USER appuser

EXPOSE 8082

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8082/health || exit 1

CMD ["./main"]