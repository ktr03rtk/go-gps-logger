FROM golang:alpine

WORKDIR /go/src/app
COPY . .

RUN cd cmd/datalogger && go build

CMD ["/bin/sh"]
