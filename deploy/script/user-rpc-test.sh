#!/bin/bash
reso_addr='crpi-qjc1bm88klc3y8ne.cn-hangzhou.personal.cr.aliyuncs.com/imdct/user-rpc-dev'
tag='latest'

container_name="im-user-rpc-test"

docker stop ${container_name}

docker rm ${container_name}

docker rmi ${reso_addr}:${tag}

docker pull ${reso_addr}:${tag}


# 如果需要指定配置文件的
# docker run -p 10001:8080 --network imooc_easy-chat -v /easy-chat/config/user-rpc:/user/conf/ --name=${container_name} -d ${reso_addr}:${tag}
docker run -p 10000:10000  --name=${container_name} --restart=always -d ${reso_addr}:${tag}
