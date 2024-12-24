goctl rpc protoc --go_out=./apps/user/rpc/ --go-grpc_out=./apps/user/rpc/ --zrpc_out=./apps/user/rpc/  ./apps/user/rpc/user.proto
goctl model mysql ddl --src ./deploy/sql/user.sql --dir ./apps/user/models/ -c
goctl api go --dir ./apps/user/api/ --style gozero -api ./apps/user/api/user.api