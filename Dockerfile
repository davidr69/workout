# Stage 1: Build the binary
ARG GO_VERSION=1.25.5
FROM docker.io/golang:${GO_VERSION}-alpine AS builder

# Install git/ca-certificates if needed for private modules or HTTPS
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy dependency files first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build
COPY . .
# CGO_ENABLED=0 creates a statically linked binary (ideal for Alpine/Scratch)
RUN CGO_ENABLED=0 GOOS=linux go build -o gym ./main.go

# Stage 2: Runtime
FROM docker.io/alpine:3.19

# Security best practice: include CA certs for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy only the compiled binary from the builder stage
COPY --from=builder /app/gym .
COPY --from=builder /app/ui ./ui

USER nobody

ENTRYPOINT ["./gym"]
