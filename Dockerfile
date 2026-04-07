FROM m.daocloud.io/docker.io/alpine:3.21

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk --no-cache add tzdata ca-certificates curl \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

ARG TARGETARCH
RUN echo "I'm building for $TARGETARCH"

COPY ./app-${TARGETARCH} /app

CMD ["/app"]
