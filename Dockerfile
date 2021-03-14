FROM golang:alpine

WORKDIR /go/src/app
COPY . .

RUN go build -o cmd/datalogger/datalogger cmd/datalogger/main.go

CMD ["/bin/sh"]
