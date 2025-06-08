# Multi-stage build برای بهینه‌سازی سایز image
FROM golang:1.24.3 AS builder

# نصب dependencies مورد نیاز
RUN apk add --no-cache git ca-certificates tzdata

# تنظیم working directory
WORKDIR /app

# کپی go mod files
COPY go.mod go.sum ./

# دانلود dependencies
RUN go mod download

# کپی source code
COPY . .

# Build application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o main .

# Final stage - استفاده از alpine برای کاهش سایز
FROM alpine:latest

# نصب ca-certificates برای HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# ایجاد non-root user برای امنیت
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /root/

# کپی binary از builder stage
COPY --from=builder /app/main .

# تغییر ownership به appuser
RUN chown appuser:appgroup main

# تغییر به non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run application
CMD ["./main"]