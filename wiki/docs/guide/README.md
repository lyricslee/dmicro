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
  -v /var/docker/nats-streaming:/datastore \
  nats-streaming:latest -store file -dir datastore
```
## database
```
import sql scripts
```
## gate micro
```
cd gate/micro
go build
./micro web
```
## web
```shell
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
```
curl -X POST \
  http://localhost:8080/xdd/passport/smslogin \
  -H 'Accept: */*' \
  -H 'App-Id: 1' \
  -H 'App-Version: 1.0.0' \
  -H 'Cache-Control: no-cache' \
  -H 'Connection: keep-alive' \
  -H 'Content-Type: application/json' \
  -H 'Device-Id: abc' \
  -H 'Host: localhost:8080' \
  -H 'Model: ' \
  -H 'Net: WIFI' \
  -H 'Os-Type: ios' \
  -H 'Os-Version: 10.0' \
  -H 'Postman-Token: 58767350-445b-411e-b6ac-98ea2989b673,95ecd2f6-ed52-4c35-9e4b-116dbe82e913' \
  -H 'Resolution: 800*600' \
  -H 'Token: ' \
  -H 'User-Agent: PostmanRuntime/7.13.0' \
  -H 'User-Id: 0' \
  -H 'accept-encoding: gzip, deflate' \
  -H 'cache-control: no-cache' \
  -d '{
	"mobile": "13803456789",
	"code": "123456"
}'
```