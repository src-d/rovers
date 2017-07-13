FROM alpine:3.6
MAINTAINER source{d}

RUN apk add --no-cache ca-certificates

ADD build/rovers_linux_amd64/rovers /bin/

CMD ["rovers","repos"]
