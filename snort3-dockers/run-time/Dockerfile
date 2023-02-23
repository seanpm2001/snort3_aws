ARG UBUNTU_VERSION="focal-20230126"
FROM ubuntu:${UBUNTU_VERSION}

ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt -y install libpcap-dev \
	libssl-dev supervisor ca-certificates

COPY --from=snort3-build /packages/ /tmp/
RUN tar -xvpf /tmp/libpcre.tar.gz -C /
RUN tar -xvpf /tmp/libluajit.tar.gz -C /
RUN tar -xvpf /tmp/libhwloc.tar.gz -C /
RUN tar -xvpf /tmp/libdnet.tar.gz -C /
RUN tar -xvpf /tmp/snort3.tar.gz -C /

RUN mkdir -p /var/log/supervisor
RUN mkdir -p /var/snort
RUN mkdir -p /var/snort/events
COPY artifact/supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY artifact/runsnort.sh /usr/local/bin/.
COPY artifact/apiagent /usr/local/bin/.

# Install lightspd template
COPY artifact/Talos_LightSPD_template.tar.gz /tmp/.
RUN tar -xvpf /tmp/Talos_LightSPD_template.tar.gz -C /var/snort/

RUN rm /tmp/*tar.gz
RUN ldconfig
