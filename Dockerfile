FROM golang:alpine

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go build -o cmd/datalogger/datalogger cmd/datalogger/main.go

CMD ["/bin/sh"]
