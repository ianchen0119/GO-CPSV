# Source: https://hub.docker.com/r/adrianszwej/opensaf
FROM ubuntu:16.04

LABEL maintainer.name="Ian Chen"
LABEL maintainer.mail="ychen.desl@gmail.com"
LABEL image.description="Image containing opensaf middleware with basic services"
LABEL container.run="NODE=SC-1 && docker run -t --name \$NODE -h \$NODE -v /home/ianchen0119/sharedfs:/etc/opensaf/sharedfs -i <img>"

RUN mkdir -p /home/opensaf
RUN mkdir -p /home/opensaf/GO-CPSV
WORKDIR /home/opensaf

# Add binaries and scripts needing for runtime
ADD https://gist.githubusercontent.com/ianchen0119/dd9aa54617c56cc57ae19962adb24176/raw/b74d6a805fe6d3b40d0a8010ab74bde2fc4e793d/setup-opensaf-node /etc/init.d/setup-opensaf-node
ADD https://raw.githubusercontent.com/ianchen0119/GO-CPSV/master/tipc-config /sbin/tipc-config

RUN chmod u+x /sbin/tipc-config \
    && chmod u+x /etc/init.d/setup-opensaf-node

# Dynamic part which can be sent via --build-arg version=5.2.GA
ARG buildversion="5.2.GA"
ARG configureflags="--enable-imm-pbe --enable-tipc"

# Packages for development. Branch "default" will be build, which is the latest. One can also use opensaf-4.6.x | opensaf-4.5.x | opensaf-4.4.x | opensaf-4.3.x instead. 
RUN apt-get update && apt-get install --no-install-recommends -y \
    sudo \
    sqlite3 \ 
    libxml2 \
    psmisc \
    python2.7-minimal \ 
    net-tools \ 
    kmod \
    ca-certificates \
    mercurial gcc g++ libxml2-dev automake m4 autoconf libtool pkg-config make python-dev libsqlite3-dev binutils git wget \
    && apt-get autoremove -y \
    && apt-get clean \ 
    && rm /var/lib/apt/lists/*.lz4 \
    && cd /home/opensaf \
    && hg clone http://hg.code.sf.net/p/opensaf/staging opensaf-staging \
    && cd opensaf-* \
    && hg update $buildversion \
    && ./bootstrap.sh \
    && ./configure $configureflags \
    && make install \
    && ldconfig \
    && sed '/\. \/lib\/lsb\/init-functions/ a\\/etc\/init.d\/setup-opensaf-node' -i /etc/init.d/opensafd

#download golang-1.17
Run wget https://dl.google.com/go/go1.17.8.linux-amd64.tar.gz && \
    tar -C /usr/local -zxvf go1.17.8.linux-amd64.tar.gz && \
    mkdir -p /go/{bin,pkg,src} && \
    rm go1.17.8.linux-amd64.tar.gz

RUN groupadd opensaf && \
 useradd -r -g opensaf -d /var/run/opensaf -s /sbin/nologin -c "OpenSAF" opensaf && \
 echo '%{opensaf_user} ALL = NOPASSWD: /sbin/reboot, /sbin/tipc-config, /usr/bin/pkill, /usr/bin/killall' >> /etc/sudoers && \
 echo 'Defaults:%opensaf !requiretty' >> /etc/sudoers && \
 echo 'Defaults:opensaf !requiretty' >> /etc/sudoers

ENV container docker
ENV GOPATH=/go
ENV GOROOT=/usr/local/go    
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH
ENV GO111MODULE=auto

COPY ./. /home/opensaf/GO-CPSV

CMD ["/bin/bash"]