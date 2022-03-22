FROM golang:latest

RUN apt update
RUN apt install -y gdal-bin
RUN go get -u github.com/derekparker/delve/cmd/dlv
