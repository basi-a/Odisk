FROM fedora:38
# 安装 HAProxy、Keepalived 和其他必要软件
RUN dnf install -y haproxy keepalived iproute procps && \
    dnf clean all 
RUN systemctl enable haproxy keepalived
# 复制启动脚本到镜像内
COPY ./scripts/check_haproxy.sh /opt/check_haproxy.sh
RUN chmod 755 /opt/check_haproxy.sh && mkdir /var/log/keepalived

# 设置启动命令
CMD ["/sbin/init"]
