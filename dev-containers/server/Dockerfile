FROM fedora:38
RUN dnf install -y iproute procps && \
    dnf clean all 
COPY ./odisk /opt/odisk
RUN chmod +x /opt/odisk && mkdir /var/log/odisk /etc/odisk

WORKDIR /root

CMD [ "/opt/odisk" ]