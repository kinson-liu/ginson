FROM golang:alpine as builder

WORKDIR /ginson
COPY . .

RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w CGO_ENABLED=0 \
    && go env \
    && go mod tidy \
    && go build -o server .

FROM alpine:latest
LABEL authors="ginson"

ENV TZ=Asia/Shanghai
RUN apk update && apk add --no-cache tzdata openntpd \
    && ln -sf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone \
WORKDIR /ginson

COPY --from=0 /ginson/server ./
COPY --from=0 /ginson/config/config.yaml ./

EXPOSE 8080
ENTRYPOINT ./server -c config.yaml