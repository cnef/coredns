FROM registry-in.dustess.com/base/alpine:3.13

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
    && apk add --update --no-cache --no-progress ca-certificates tzdata \
    && update-ca-certificates \
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && rm -rf /var/cache/apk/*

ADD coredns /coredns

EXPOSE 53 53/udp
ENTRYPOINT ["/coredns"]
