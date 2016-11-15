# vpnkeeper
Mac下基于命令行的VPN拨号工具，解决下载时间过长，VPN断线导致下载中断的问题。

## 编译过程

### 1. Go语言环境

```bash
# 以root用户身份登录
# 下载
$ wget https://storage.googleapis.com/golang/go1.7.1.linux-amd64.tar.gz
$ tar xzvf go1.7.1.linux-amd64.tar.gz
$ mv go /usr/local

# 设置GO环境
$ vi ~/.bash_profile

# 添加如下代码
export GOROOT=/usr/local/go
export GOPATH=$HOME/go-project
export PATH=$PATH:$HOME/bin:$GOPATH/bin

# 生效GO PATH配置
$ source ~/.bash_profile
```

### 2. Godep安装

Godep是一个go语言库包管理工具。

```bash
$ go get github.com/tools/godep
```

### 3. 下载工程

* 把项目路径加入到$GOPATH/src
* 依赖的项目和项目本身都应该是个git仓库
* 目录结构例如

```
$GOPATH
 |-src
 |  |-vpnkeeper
 |-pkg
 |-bin

```

```bash
$ cd $GOPATH/src
$ git clone https://github.com/lordking/vpnkeeper
```

### 4. 编译

```bash
$ cd vpnkeeper
$ godep restore
$ godep go build
```
