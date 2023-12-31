FROM golang:1.20.4-alpine3.18 as builder
COPY go.mod go.sum /ftp_driver/
WORKDIR /ftp_driver
RUN go mod download
COPY . /ftp_driver

RUN go build -o driver ./cmd/driver/main.go

FROM alpine

RUN echo "@testing https://dl-cdn.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories && \
	apk update && \
    apk add curlftpfs@testing

RUN mkdir -p /var/run/docker/ftp-driver/ /var/run/docker/ftp-driver/state

WORKDIR /

COPY --from=builder /ftp_driver/driver .
CMD ["driver"]
