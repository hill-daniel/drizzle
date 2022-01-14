# golang:1.17-alpine
FROM golang@sha256:8c7994cfdf4d488799d40d85d83bd41c7fd290e8eed1affc2abd386150750d2d as builder

RUN apk update

WORKDIR /app

COPY . .

RUN go mod download
RUN go mod verify

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" ./cmd/drizzle

# Build image
FROM golang@sha256:8c7994cfdf4d488799d40d85d83bd41c7fd290e8eed1affc2abd386150750d2d
RUN apk update  \
    && apk add --no-cache \
           git \
           build-base \
           hugo \
           python3 \
           py3-pip \
    && pip3 install --upgrade pip \
    && pip3 install --no-cache-dir \
            awscli \
    && rm -rf /var/cache/apk/*
RUN git config --global user.email "drizzle@uphill.dev"
RUN git config --global user.name "drizzle"
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

COPY --from=builder /app/drizzle /app/

ENTRYPOINT ["/app/drizzle"]
