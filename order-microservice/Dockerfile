# Stage 1: Build
FROM golang:1.24-alpine AS builder


WORKDIR /app

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the full source including .env
COPY . .

# Build the Go app
RUN go build -o server .

# Stage 2: Run
FROM alpine:latest

WORKDIR /root

# Copy built binary and .env file
COPY --from=builder /app/server /server
COPY --from=builder /app/.env .env 

COPY --from=builder /app/go-microservice-rohit.json /app/credentials.json

# Set Google credentials ENV variable
ENV GOOGLE_APPLICATION_CREDENTIALS=/app/credentials.json

EXPOSE 8080
CMD ["/server"]