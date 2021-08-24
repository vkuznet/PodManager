# Start from the latest golang base image
FROM golang:latest as builder
MAINTAINER Valentin Kuznetsov vkuznet@gmail.com
ENV WDIR=/data
WORKDIR $WDIR

# Install latest kubectl for using with crons
RUN curl -ksLO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl
RUN chmod +x kubectl

# RUN go get github.com/vkuznet/PodManager
ARG CGO_ENABLED=0
RUN git clone https://github.com/vkuznet/PodManager.git && cd PodManager && make

# FROM alpine
# RUN mkdir -p /data
# https://blog.baeke.info/2021/03/28/distroless-or-scratch-for-go-apps/
FROM gcr.io/distroless/static AS final
COPY --from=builder /data/PodManager/PodManager /data/
COPY --from=builder /data/kubectl /usr/bin/
