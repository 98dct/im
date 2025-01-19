#!/bin/bash
reso_addr='crpi-qjc1bm88klc3y8ne.cn-hangzhou.personal.cr.aliyuncs.com/imdct/im-ws-dev'
tag='latest'

POD_IP="192.168.8.100"

container_name="im-im-ws-test"

docker stop ${container_name}

docker rm ${container_name}

docker rmi ${reso_addr}:${tag}

docker pull ${reso_addr}:${tag}


# 如果需要指定配置文件的
# docker run -p 10001:8080 --network imooc_easy-chat -v /easy-chat/config/user-rpc:/user/conf/ --name=${container_name} -d ${reso_addr}:${tag}
docker run -p 10090:10090 -e POD_IP=${POD_IP} --name=${container_name} --restart=always -d ${reso_addr}:${tag}