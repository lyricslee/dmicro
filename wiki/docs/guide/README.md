# 快速开始

## jaeger
```
docker run -d --name jaeger \
  -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
  -p 5775:5775/udp \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 14268:14268 \
  -p 9411:9411 \
  jaegertracing/all-in-one:latest
```
## nats-streaming
```
docker run -d --name nats-streaming \
  -p 4222:4222 \
  -p 8222:8222 \
  -v /var/lib/docker-data/nats-streaming:/datastore \
  nats-streaming:latest -store file -dir datastore
```
## etcd
```
docker run -d --name etcd \
  -p 2379:2379 \
  -v /var/lib/docker-data/etcd:/etcd-data \
  quay.io/coreos/etcd:latest \
  etcd \
  --data-dir=/etcd-data \
  --name node1 \
  --listen-client-urls http://0.0.0.0:2379 \
  --advertise-client-urls http://0.0.0.0:2379
```
## database
```
import sql scripts
```
## gate micro web
```
cd gate/micro
go build
./micro web
```
## web
```
cd web/dd
go build
./dd
```
## gid
```
cd srv/gid
go build
./gid
```
## passport
```
cd srv/passport
go build
./passport
```
## user
```
cd srv/user
go build
./user
```
## postman
接口地址注意大小写，默认端口8082
```
curl -X POST \
  'http://localhost:8082/dd/passport/SmsLogin' \
  -H 'Content-Type: application/json' \
  -H 'App-Id: 1' \
  -d '{
    "mobile": "13705918888",
    "code": "123456"
}'
```
以上使用的是web服务作为聚合api接口，还可以使用api服务作为聚合api接口

## gate micro api
```
./micro api --handler=api
```

## api
```
cd api/dd
go build
./dd
```

## postman
接口地址注意大小写，默认端口8080
```
curl -X POST \
  'http://localhost:8080/dd/passport/SmsLogin' \
  -H 'Content-Type: application/json' \
  -H 'App-Id: 1' \
  -d '{
    "mobile": "13705918888",
    "code": "123456"
}'
```