FROM golang:1.16.5-alpine3.14 AS builder
WORKDIR /walletservice
COPY . .
RUN CGO_ENABLED=0 go build \
    -o /bin/walletservice \
    /walletservice/cmd/walletservice/

FROM alpine:3.14
COPY --from=builder /bin/walletservice /opt
CMD ["/opt/walletservice"]
