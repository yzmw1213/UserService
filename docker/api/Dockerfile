FROM golang:1.12

ENV PKG_PATH /go/src/app
ENV SERVICE_PATH $PKG_PATH/project
ENV GO111MODULE on
ENV DOCKERIZE_VERSION v0.6.1
WORKDIR /go/src/app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . $PKG_PATH

RUN apt-get update && apt-get install -y wget
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz
