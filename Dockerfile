# ---- Build Stage ----
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install CA certs and create non-root user
RUN apk add --no-cache ca-certificates git && \
    addgroup -S appgroup && \
    adduser -S -D -G appgroup -s /bin/sh -u 1001 appuser

# Copy go.mod first for better caching
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build statically-linked binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -a -installsuffix cgo \
    -o backend main.go

# ---- Final Stage ----
FROM scratch AS final

WORKDIR /

# Copy essential files from builder
COPY --from=builder /app/backend /backend
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Use non-root user
USER 1001

# Expose port
EXPOSE 8080

# Run the binary
ENTRYPOINT ["/backend"]