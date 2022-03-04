FROM golang:latest

RUN apt update
RUN apt install -y gdal-bin
