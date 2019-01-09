FROM golang:1.11.4-alpine

WORKDIR /src/github.com/LevInterctive/imageup
ADD . /src/github.com/LevInterctive/imageup

RUN apk add bash ca-certificates git gcc g++ libc-dev

RUN go build -o bin/imageup
CMD ["./bin/imageup"]
