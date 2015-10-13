FROM tyba/base

MAINTAINER Tyba

RUN echo Europe/Madrid | tee /etc/timezone && dpkg-reconfigure --frontend noninteractive tzdata

WORKDIR /opt

RUN apt-get update -y && apt-get install wget bzip2 -y
RUN wget https://github.com/aktau/github-release/releases/download/v0.5.3/linux-amd64-github-release.tar.bz2
RUN bunzip2 linux-amd64-github-release.tar.bz2 && tar -xf linux-amd64-github-release.tar

ENV ENVIRONMENT=production
ENV ETCD_SERVERS=http://etcd.oss.tyba.cc:2379
ENV GITHUB_TOKEN=08763897c930b3ff7f7cebf8da45935350a96b7d
ENV GITHUB_USER=tyba
ENV GITHUB_REPO=srcd-rovers

RUN /opt/bin/linux/amd64/github-release download -t build -n srcd-rovers_linux_amd64.tar.gz && echo $DOCKERSHIP_REV
RUN tar -xvzf srcd-rovers_linux_amd64.tar.gz && \
	rm -f srcd-rovers_linux_amd64.tar.gz && \
	chown root:root -R srcd-rovers_linux_amd64

CMD ["bash", "-c"]
