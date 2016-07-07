FROM quay.io/srcd/basic:latest
MAINTAINER source{d}

ADD bin /bin

ENV ENVIRONMENT=production
ENV ETCD_SERVERS=http://etcd.oss.tyba.cc:4001

WORKDIR /opt

#RUN apt-get install -y wget \
#  && wget https://github.com/mcuadros/ofelia/releases/download/v0.2.1/ofelia_v0.2.1_linux_amd64.tar.gz -O ofelia.tar.gz && tar -xvzf ofelia.tar.gz \
#  && rm -rf /var/lib/apt/lists/*
#
#ADD ofelia.ini /etc/ofelia.ini
#VOLUME /var/log/sync/

#CMD ["/opt/ofelia_linux_amd64/ofelia", "daemon", "--config", "/etc/ofelia.ini"]
