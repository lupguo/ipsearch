# Golang builder service
FROM golang:1.12-alpine as go-builder
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk add build-base autoconf \
    && apk add git curl tcpdump
COPY . /data/go/ip-search/
WORKDIR /data/go/ip-search/
RUN export GOPROXY="https://goproxy.io" \
    && export CGO_ENABLED="0" \
    && go get -u -v ...
#RUN export GOPROXY="https://goproxy.io" \
#    && export CGO_ENABLED="0" \
#    && go get -u -v github.com/tkstorm/ip-serach/...

# Running application service
FROM alpine:latest as goapp
COPY --from=go-builder /go/bin/ /go/bin/
EXPOSE 8680
CMD ["/go/bin/ipshttpd", "-listen", ":8680"]
