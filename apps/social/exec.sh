goctl rpc protoc --go_out=./apps/social/rpc/ --go-grpc_out=./apps/social/rpc/ --zrpc_out=./apps/social/rpc/  ./apps/social/rpc/social.proto
goctl model mysql ddl --src ./deploy/sql/social.sql --dir ./apps/social/socialmodels/ -c
goctl api go --dir ./apps/social/api/ --style gozero -api ./apps/social/api/social.api