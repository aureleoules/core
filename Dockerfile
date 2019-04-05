FROM golang:alpine AS binaryBuilder
# Install build deps
ADD https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep
RUN apk add --no-cache git
# Build project
WORKDIR /go/src/github.com/backpulse/core
COPY . .
RUN dep ensure
RUN go build -o main .

FROM alpine:latest
# Copy target app from binaryBuilder stage
WORKDIR /app/backpulse
COPY --from=binaryBuilder /go/src/github.com/backpulse/core/main .
# Configure Docker Container
EXPOSE 8000
CMD ["/app/backpulse/main"]