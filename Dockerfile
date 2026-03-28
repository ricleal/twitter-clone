# syntax=docker/dockerfile:1

###############
# Build stage #
###############
FROM golang:1.25.0-alpine AS builder

WORKDIR /app

# Download dependencies separately for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY cmd/ cmd/
COPY internal/ internal/

# Build static binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /app/app ./cmd/twitter



#################
# Runtime stage #
#################
FROM alpine:3.20

ARG API_PORT=8888
ENV API_PORT=${API_PORT}

# Install minimal runtime dependencies
RUN apk add --no-cache curl ca-certificates

# Create non-root user
RUN addgroup -g 1000 app && adduser -u 1000 -G app -s /sbin/nologin -D app

COPY --from=builder --chmod=755 /app/app /app/app

EXPOSE ${API_PORT}

USER app

HEALTHCHECK --interval=5s --timeout=3s --retries=3 \
    CMD curl -f http://localhost:${API_PORT}/health || exit 1

ENTRYPOINT ["sh", "-c", "exec /app/app -port \"${API_PORT}\" \"$@\"", "--"]
