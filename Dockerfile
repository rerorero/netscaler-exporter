FROM golang:1.9.2 as builder

RUN wget -O /usr/bin/dep https://github.com/golang/dep/releases/download/v0.3.2/dep-linux-amd64 && \
    echo '322152b8b50b26e5e3a7f6ebaeb75d9c11a747e64bbfd0d8bb1f4d89a031c2b5  /usr/bin/dep' | sha256sum -c - && \
    chmod +x /usr/bin/dep

RUN mkdir -p /go/src/github.com/rerorero/netscaler-exporter
WORKDIR /go/src/github.com/rerorero/netscaler-exporter

COPY Gopkg.toml Gopkg.lock ./

RUN dep ensure -vendor-only
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -a -installsuffix cgo -o netscaler-exporter

FROM alpine:latest
EXPOSE 8080

RUN apk --no-cache add ca-certificates && \
    mkdir -p /etc/nsx/

WORKDIR /root/

COPY --from=builder /go/src/github.com/rerorero/netscaler-exporter/netscaler-exporter .
CMD ./netscaler-exporter --conf.file=/etc/nsx/nsxconf.yml
