# 一个web笔记本

本站demo：长期维护 https://note.valarx.com

## 更新记录

2021/09/10

添加通过API POST请求创建

```
参数
Tag 长度控制请修改 Api() 中的值，默认不小于10
Text 默认不为空即可

```

原项目地址：https://github.com/pereorga/minimalist-web-notepad

## 使用说明:

go version 1.16

修改代码中func NewDBEngine()中数据库相关配置

export GO111MODULE=on

export GOPROXY=https://goproxy.cn

默认端口23456

go build后将文件夹二进制要放在note目录下

nohup ./note &

