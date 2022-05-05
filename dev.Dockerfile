FROM golang:latest

RUN apt update
RUN apt -y install gdal-bin gdal-data libgdal-dev
