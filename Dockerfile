FROM golang:1.11-alpine

ADD . /go/src/github.com/LevInterctive/imageup/
WORKDIR /go/src/github.com/LevInterctive/imageup/

RUN apk add --no-cache git openssl bzr \
    && go get -u github.com/golang/dep/cmd/dep

RUN dep ensure && go build -o bin/imageup
CMD ["./bin/imageup"]
