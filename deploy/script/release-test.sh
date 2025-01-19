need_start_server=(
  user-rpc-test.sh
  user-api-test.sh
  social-rpc-test.sh
  social-api-test.sh

  im-api-test.sh
  im-rpc-test.sh
  im-ws-test.sh
  task-mq-test.sh
)

for i in ${need_start_server[*]} ; do
    chmod +x $i
    ./$i
done

docker ps
docker exec -it etcd_3_5_7 etcdctl get "" --prefix