# vagrant环境搭建

## VirtualBox

### 安装VirtualBox
下载地址：
https://www.virtualbox.org/

<img src="https://github.com/fztcjjl/dmicro/raw/master/docs/drone/img/vagrant/1.png">

### 设置默认虚拟机位置

<img src="https://github.com/fztcjjl/dmicro/raw/master/docs/drone/img/vagrant/2.png">

## vagrant

### 安装Vagrant

https://www.vagrantup.com/

### 添加Box

CentOS 7 box地址：

https://app.vagrantup.com/centos/boxes/7/versions/1902.01

<img src="https://github.com/fztcjjl/dmicro/raw/master/docs/drone/img/vagrant/3.png">

由于在vagrant box hub上下载box速度很慢，这里提前下载好centos 7 boxes，**保存于E:\vagrant\box_files文件夹下，注意：目录结构保持与我的一致。**
**下载方法是box地址+提供商名字**
**如：https://app.vagrantup.com/centos/boxes/7/versions/1902.01/providers/virtualbox.box**

<img src="https://github.com/fztcjjl/dmicro/raw/master/docs/drone/img/vagrant/4.png">

virtualbox_centos-7.4.json文件内容如下：

```json
{
  "name": "centos/7",
  "versions": [
    {
      "version": "1902.01",
      "providers": [
        {
          "name": "virtualbox",
          "url": "file:///E:/vagrant/box_files/CentOS-7-x86_64-Vagrant-1902_01.VirtualBox.box"
        }
      ]
    }
  ]
}
```

在命令行下进入E:\vagrant目录执行
**vagrant.exe box add box_files\CentOS-7-x86_64-Vagrant-1902_01.VirtualBox.json**

<img src="https://github.com/fztcjjl/dmicro/raw/master/docs/drone/img/vagrant/5.png">

### 通过vagrant启动虚拟机

进入lab-metoo文件夹，准备好Vagrantfile与setup.sh文件，运行vagrant.exe up启动虚拟机。

<img src="https://github.com/fztcjjl/dmicro/raw/master/docs/drone/img/vagrant/6.png">

Vagrntfile文件内容如下：

```
# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.require_version ">= 1.6.0"

boxes = [
    {
        :name => "metoo",
        :eth1 => "192.168.56.186",
        :mem => "2048",
        :cpu => "2"
    }
]

Vagrant.configure(2) do |config|
  config.vm.box = "centos/7"
  config.vbguest.auto_update = false
  config.vbguest.no_remote = true

  boxes.each do |opts|
    config.vm.define opts[:name] do |config|
      config.vm.hostname = opts[:name]
      config.vm.provider "vmware_fusion" do |v|
        v.vmx["memsize"] = opts[:mem]
        v.vmx["numvcpus"] = opts[:cpu]
      end

      config.vm.provider "virtualbox" do |vb|
        vb.name = "lab-#{config.vm.hostname}"
        vb.customize ["modifyvm", :id, "--memory", opts[:mem]]
        vb.customize ["modifyvm", :id, "--cpus", opts[:cpu]]
      end

      config.vm.provider :libvirt do |lv|
          lv.host = "lab-#{config.vm.hostname}"
      end

      config.vm.network :private_network, ip: opts[:eth1]
    end
  end

  # 禁用vagrant的默认共享目录
  # config.vm.synced_folder ".", "/vagrant", disabled:true
  # config.vm.synced_folder "../public", "/opt", mount_options:["dmode=777","fmode=666"]
  config.vm.provision "shell", privileged: true, path: "./setup.sh"
end
```

setup.sh文件内容如下:
```
#/bin/sh

#  Delta RPMs disabled because /usr/bin/applydeltarpm not installed.
sudo yum install -y deltarpm

# 设置时区
timedatectl set-timezone Asia/Shanghai

# 关闭防护墙及selinux
sed -i '/SELINUX/s/enforcing/disabled/g' /etc/selinux/config
setenforce 0
systemctl stop firewalld.service
systemctl disable firewalld.service

# 允许密码登录
sed -i '/^PasswordAuthentication no/s/no/yes/g' /etc/ssh/sshd_config
systemctl restart sshd.service

# install some tools
sudo yum install -y git vim gcc glibc-static telnet bridge-utils net-tools

# install docker
## step 1: 安装必要的一些系统工具
sudo yum install -y yum-utils device-mapper-persistent-data lvm2
## Step 2: 添加软件源信息
sudo yum-config-manager --add-repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo
## Step 3: 更新并安装 Docker-CE
sudo yum makecache fast
sudo yum -y install docker-ce
## 镜像加速器
sudo mkdir -p /etc/docker
sudo tee /etc/docker/daemon.json <<-'EOF'
{
  "registry-mirrors": ["https://c8x7u9xt.mirror.aliyuncs.com"]
}
EOF

# start docker service
sudo usermod -aG docker vagrant
sudo systemctl start docker
sudo systemctl enable docker

# install docker-compose
curl -L https://get.daocloud.io/docker/compose/releases/download/1.25.1/docker-compose-`uname -s`-`uname -m` > /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

```
### 安装vagrant-vbguest插件

进入命令行执行：
```
vagrant.exe plugin install vagrant-vbguest
vagrant.exe vbguest
```
安装vagrant-vbguest插件用于文件夹共享
**如果不需要文件夹共享，可以省略，这个安装过程也比较慢，如果安装时间过长，建议直接跳过这一步**

### 通过vagrant ssh登录虚拟机
虚拟机安装好后，默认已经创建了vagrant用户，密码是vagrant，root用户的默认密码也是vagrant。
在命令行下执行**vagrant.exe ssh**，将以vagrant用户登录虚拟机。

<img src="https://github.com/fztcjjl/dmicro/raw/master/docs/drone/img/vagrant/7.png">

也可以用Secure CRT、xshell等工具连接登录。

### 常见命令

| **命令**           | **说明**              |
| ------------------ | --------------------- |
| vagrant box list   | 查看目前已有的box     |
| vagrant box add    | 新增加一个box         |
| vagrant box remove | 删除指定box           |
| vagrant init       | 初始化配置Vagrantfile |
| vagrant up         | 启动虚拟机            |
| vagrant ssh        | ssh登录虚拟机         |
| vagrant suspend    | 挂起虚拟机            |
| vagrant reload     | 重启虚拟机            |
| vagrant halt       | 关闭虚拟机            |
| vagrant status     | 查看虚拟机状态        |
| vagrant destroy    | 删除虚拟机            |
