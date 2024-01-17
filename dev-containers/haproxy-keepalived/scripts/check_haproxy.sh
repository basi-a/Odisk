#!/bin/bash
echo "$(date) start check haproxy..." >> /var/log/keepalived/check_haproxy.log
HA=$(ps -C haproxy --no-header | wc -l)
if [ "$HA" -eq 0 ];then
        echo "$(date) haproxy is not running, restarting..." >> /var/log/keepalived/check_haproxy.log
        systemctl start haproxy
        if [ "$HA" -eq 0 ];then
            echo "$(date) haproxy is down, kill all keepalived..." >> /var/log/keepalived/check_haproxy.log
            killall keepalived
        fi
fi
