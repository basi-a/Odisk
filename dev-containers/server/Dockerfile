FROM fedora:38
 
COPY ./odisk /opt/odisk

RUN chmod +x /opt/odisk && mkdir /var/log/odisk /etc/odisk /opt/cert /opt/template
COPY ./cert/* /opt/cert/
COPY ./template/* /opt/template/
WORKDIR /root

CMD [ "/opt/odisk" ]
