# ==============================
# 1) BUILD STAGE
# ==============================
ARG GO_VERSION=1.25.7
FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /service
COPY . .

RUN go mod vendor && \
    CGO_ENABLED=0 GOOS=linux go build \
    -mod=vendor \
    -o main \
    ./cmd/uni-products-svc

# ==============================
# 2) FINAL STAGE
# ==============================
FROM alpine:latest

WORKDIR /service

RUN apk add --no-cache ca-certificates

COPY --from=builder /service/main .
COPY --from=builder /service/config.yaml .

ENV KV_VIPER_FILE=/service/config.yaml

CMD ["./main", "run", "service"]

