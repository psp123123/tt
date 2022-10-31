#!/bin/bash
echo 'vm.overcommit_memory = 1' >> /etc/sysctl.conf
echo 'net.core.somaxconn=65535' >> /etc/sysctl.conf
sysctl -p
/usr/sbin/useradd redis -s /sbin/nologin
chown -R redis:redis /home/software/redis
echo 'export PATH=$PATH:/home/software/redis/bin/' >> /etc/profile
. /etc/profile
mv /home/software/redis6379.service /etc/systemd/system/