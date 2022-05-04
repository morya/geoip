FROM alpine:3.11
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

COPY geoip /app
COPY ip-database /ip-database


CMD [ "/app" ]
