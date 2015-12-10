FROM tyba/base

MAINTAINER Tyba

WORKDIR /opt

RUN apt-get update -y && apt-get install wget bzip2 -y
RUN wget https://github.com/aktau/github-release/releases/download/v0.5.3/linux-amd64-github-release.tar.bz2
RUN bunzip2 linux-amd64-github-release.tar.bz2 && tar -xf linux-amd64-github-release.tar

ENV ENVIRONMENT=production
ENV ETCD_SERVERS=http://etcd.oss.tyba.cc:4001
ENV GITHUB_TOKEN=08763897c930b3ff7f7cebf8da45935350a96b7d
ENV GITHUB_USER=src-d
ENV GITHUB_REPO=rovers

RUN TAG=`/opt/bin/linux/amd64/github-release info | sed -n 2p | cut -d " " -f 2` \
    && /opt/bin/linux/amd64/github-release download -t $TAG -n rovers_${TAG}_linux_amd64.tar.gz \
    && echo $DOCKERSHIP_REV

RUN tar -xvzf rovers_v*_linux_amd64.tar.gz && \
	rm -f rovers_v*_linux_amd64.tar.gz && \
	chown root:root -R rovers_linux_amd64

CMD ["bash", "-c"]
