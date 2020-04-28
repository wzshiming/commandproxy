FROM golang:alpine AS builder
WORKDIR /go/src/github.com/wzshiming/commandproxy/
COPY . .
ENV CGO_ENABLED=0
RUN go install ./cmd/commandproxy

FROM alpine
COPY --from=builder /go/bin/commandproxy /usr/local/bin/
ENTRYPOINT [ "/usr/local/bin/commandproxy" ]