FROM alpine:3.4

LABEL maintainer "beedemo.sa@gmail.com"

RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

EXPOSE 8080
ENV DB db
CMD ["go-demo"]
HEALTHCHECK --interval=10s CMD wget -qO- localhost:8080/demo/hello

ARG COMMIT_SHA
LABEL beedemo.commit.sha=$COMMIT_SHA

ARG BUILD_CACHE_COMMIT_SHA
#the beedemo.build.cache.commit.sha will allows us to recreate the build environment for this image
LABEL beedemo.build.cache.commit.sha=$BUILD_CACHE_COMMIT_SHA

COPY go-demo /usr/local/bin/go-demo
RUN chmod +x /usr/local/bin/go-demo