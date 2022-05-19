# =========================================
# build stage
# =========================================
FROM golang:1.18-alpine as build-stage

WORKDIR /opt/gps_receiver

ARG TZ
RUN apk --no-cache add tzdata && \
  cp /usr/share/zoneinfo/${TZ} /etc/localtime && \
  echo "${TZ}" > /etc/timezone && \
  apk del tzdata

COPY go.mod go.sum ./
RUN go mod download \
  && go install github.com/rakyll/gotest@latest \
  && go install github.com/cespare/reflex@latest \
  && go install github.com/golang/mock/mockgen@latest \
  && apk --no-cache add make

COPY ./ ./
RUN CGO_ENABLED=0 go build -ldflags="-extldflags=-static" -gcflags="all=-N -l" -o receiver


CMD ["/bin/sh"]

# =========================================
# production stage
# =========================================
FROM alpine as production-stage

WORKDIR /opt/gps_receiver

ARG TZ
RUN apk --no-cache add tzdata && \
  cp /usr/share/zoneinfo/${TZ} /etc/localtime && \
  echo "${TZ}" > /etc/timezone && \
  apk del tzdata

COPY --from=build-stage /opt/gps_receiver/receiver /usr/bin/receiver

ENTRYPOINT ["receiver"]
