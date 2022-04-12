FROM golang:latest

RUN apt update
RUN apt install -y gdal-bin
RUN go install github.com/go-delve/delve/cmd/dlv@latest
