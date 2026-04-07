FROM alpine:3.21

# 设置时区和中文镜像源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk --no-cache add tzdata ca-certificates curl \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

# 构建参数
ARG TARGETARCH
ARG GIT_REV
ARG BUILD_DATE

# 标签信息
LABEL org.opencontainers.image.title="GeoIP Service" \
      org.opencontainers.image.description="IP geolocation service using MaxMind GeoLite2 database" \
      org.opencontainers.image.url="https://github.com/${{ github.repository }}" \
      org.opencontainers.image.source="https://github.com/${{ github.repository }}" \
      org.opencontainers.image.version="${{ github.ref_name }}" \
      org.opencontainers.image.revision="${GIT_REV}" \
      org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.licenses="MIT"

# 复制应用程序
COPY ./app-${TARGETARCH} /app

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# 运行应用程序
USER nobody:nobody
CMD ["/app"]
