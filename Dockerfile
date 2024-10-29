# Build stage using the Go image
ARG GO_VERSION=1.22
FROM golang:${GO_VERSION} as builder

# Set up the build environment
WORKDIR /app
COPY . .

# Download dependencies
RUN go mod download

# Build the Go application binary

RUN CGO_ENABLED=0 GOARCH=amd64 go build -o /app/service ./
RUN chmod +x /app/service  # Ensure binary has execute permissions

# Final scratch-based stage
FROM scratch

# Copy required system files from the builder stage
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy only the application binary
COPY --from=builder /app/service /app/service
WORKDIR /app
CMD ["/app/service"]
