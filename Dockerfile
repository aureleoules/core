FROM golang:alpine

ADD https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

RUN apk add --no-cache git

RUN mkdir /go/src/github.com/
RUN mkdir /go/src/github.com/aureleoules
RUN mkdir /go/src/github.com/backpulse/core

ADD . /go/src/github.com/backpulse/core

WORKDIR /go/src/github.com/backpulse/core 

COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only

RUN go build -o main .

CMD ["./main"]