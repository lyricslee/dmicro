# drone安装配置

## 配置gitea OAuth2授权drone登录

<img src="https://github.com/fztcjjl/dmicro/raw/master/docs/drone/img/drone/1.png">

<img src="https://github.com/fztcjjl/dmicro/raw/master/docs/drone/img/drone/2.png">


请记住客户端ID与客户端密钥，将用于docker-compse.yml中

## docker-compose.yml

```yaml
version: '3'

services:
  mysql:
    image: mysql:5.7
    container_name: drone_mysql
    restart: always
    ports:
      - "13306:3306"
    environment:
      - TZ=Asia/Shanghai
      - MYSQL_ROOT_PASSWORD=123456
      - MYSQL_USER=drone
      - MYSQL_PASSWORD=drone
      - MYSQL_DATABASE=drone
      - MYSQL_DATABASE_CHARSET=utf8mb4
      - MYSQL_DATABASE_COLLATION=utf8mb4_general_ci
    volumes:
      - /root/docker/data/drone/mysql:/var/lib/mysql
  drone-server:
    image: drone/drone:latest
    container_name: drone
    restart: always
    ports:
      - 10080:80
      - 8443:443
    environment:
      - TZ=Asia/Shanghai
      - DRONE_AGENTS_ENABLED=true
      - DRONE_GITEA_SERVER=http://192.168.56.190:10080/
      - DRONE_GITEA_CLIENT_ID=9fee9dbb-527f-47e6-8558-8d2bee53a4b6
      - DRONE_GITEA_CLIENT_SECRET=q-3KC6ui7CRNkgdeRTbA6-5LvjjjL01i6YelhyMVj-U=
      - DRONE_RPC_SECRET=5d1789d5aa2ee55e6a5b956bec3c328f
      - DRONE_SERVER_HOST=192.168.56.189:10080
      - DRONE_SERVER_PROTO=http
      - DRONE_USER_CREATE=username:fztcjjl,admin:true
#      - DRONE_LOGS_DEBUG=true
      - DRONE_GIT_USERNAME=fztcjjl
      - DRONE_GIT_PASSWORD=123456
      - DRONE_GIT_ALWAYS_AUTH=false
      - DRONE_DATABASE_DRIVER=mysql
      - DRONE_DATABASE_DATASOURCE=drone:drone@tcp(drone_mysql:3306)/drone?parseTime=true
    volumes:
      - /root/docker/data/drone/drone:/data
    depends_on:
      - mysql
  drone-agent:
    image: drone/agent:latest
    container_name: runner
    restart: always
    ports:
      - 13000:3000
    environment:
      - TZ=Asia/Shanghai
      - DRONE_RPC_PROTO=http
      - DRONE_RPC_HOST=drone
      - DRONE_RPC_SECRET=5d1789d5aa2ee55e6a5b956bec3c328f
      - DRONE_RUNNER_CAPACITY=2
      - DRONE_RUNNER_NAME=${HOSTNAME}
      - DRONE_DEBUG=true
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    depends_on:
      - drone-server
```
## 执行docker-compose up -d

docker-compose up -d

## 配置


<img src="https://github.com/fztcjjl/dmicro/raw/master/docs/drone/img/drone/3.png">

<img src="https://github.com/fztcjjl/dmicro/raw/master/docs/drone/img/drone/4.png">

<img src="https://github.com/fztcjjl/dmicro/raw/master/docs/drone/img/drone/5.png">

<img src="https://github.com/fztcjjl/dmicro/raw/master/docs/drone/img/drone/6.png">

## .drone.yml

详见dmicro/.drone.yml

## 提交代码触发CI CD

<img src="https://github.com/fztcjjl/dmicro/raw/master/docs/drone/img/drone/7.png">
