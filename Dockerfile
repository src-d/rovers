FROM alpine:3.6
MAINTAINER source{d}

ADD build/rovers_linux_amd64/rovers /bin/

CMD ["rovers","repos"]
