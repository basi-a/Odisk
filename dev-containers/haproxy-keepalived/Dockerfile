FROM fedora:38
# 安装 HAProxy、Keepalived 和其他必要软件
RUN sed -e 's|^metalink=|#metalink=|g' \
    -e 's|^#baseurl=http://download.example/pub/fedora/linux|baseurl=https://mirrors.aliyun.com/fedora|g' \
    -i.bak \
    /etc/yum.repos.d/fedora.repo \
    /etc/yum.repos.d/fedora-modular.repo \
    /etc/yum.repos.d/fedora-updates.repo \
    /etc/yum.repos.d/fedora-updates-modular.repo
RUN dnf makecache && \
    dnf install -y haproxy keepalived iproute procps && \
    dnf clean all 
RUN systemctl enable haproxy keepalived
# 复制启动脚本到镜像内
COPY ./scripts/check_haproxy.sh /opt/check_haproxy.sh
RUN chmod 755 /opt/check_haproxy.sh && mkdir /var/log/keepalived /var/log/haproxy /opt/cert
COPY ./cert/* /opt/cert/
WORKDIR /root
# 设置启动命令
CMD ["/sbin/init"]
