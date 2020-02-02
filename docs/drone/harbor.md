# 私有仓库harbor安装配置

## 下载安装包

harbor下载地址：
https://github.com/goharbor/harbor/releases

<img src="https://github.com/fztcjjl/dmicro/raw/master/docs/gitea/drone/img/harbor/1.png">

可以选择在线版与离线版，这里我选择离线版
下载，解压：
<img src="https://github.com/fztcjjl/dmicro/raw/master/docs/gitea/drone/img/harbor/2.png">

## 编辑harbor.yml文件

<img src="https://github.com/fztcjjl/dmicro/raw/master/docs/gitea/drone/img/harbor/3.png">

将主机名称改为：**192.168.56.190**

<img src="https://github.com/fztcjjl/dmicro/raw/master/docs/gitea/drone/img/harbor/4.png">

## 安装

执行install.sh

<img src="https://github.com/fztcjjl/dmicro/raw/master/docs/gitea/drone/img/harbor/5.png">

## 登录

### docker客户端登录

docker login 192.168.56.190

<img src="https://github.com/fztcjjl/dmicro/raw/master/docs/gitea/drone/img/harbor/6.png">

**用户名: admin**

**密码: Harbor12345**

编辑 /etc/docker/daemon.json

<img src="https://github.com/fztcjjl/dmicro/raw/master/docs/gitea/drone/img/harbor/7.png">

重启docker

systemctl restart docker

再次登录

<img src="https://github.com/fztcjjl/dmicro/raw/master/docs/gitea/drone/img/harbor/8.png">

### 浏览器登录

http://192.168.56.190
