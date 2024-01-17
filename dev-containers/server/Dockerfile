FROM fedora:38

COPY ./odisk /opt/odisk
RUN chmod +x /opt/odisk && mkdir /var/log/odisk /etc/odisk

WORKDIR /root

CMD [ "/opt/odisk" ]