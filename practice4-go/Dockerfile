# Multi-stage build (small final image)

FROM golang:1.22-alpine AS builder

WORKDIR /app

# Download deps first (better caching)
COPY go.mod ./
RUN go mod download

# Copy project files
COPY . .

# Build a static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o myapp ./cmd/api


FROM alpine:3.20

RUN apk add --no-cache ca-certificates

# Copy binary + required runtime files
COPY --from=builder /app/myapp /myapp
COPY --from=builder /app/docs /docs
COPY --from=builder /app/database /database

EXPOSE 8080

ENTRYPOINT ["/myapp"]
