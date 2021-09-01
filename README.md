#  一个web笔记本
原始项目地址：https://github.com/pereorga/minimalist-web-notepad


本站地址 https://low.im
##使用说明:
go version 1.16

修改代码中func NewDBEngine()中数据库相关配置

export GO111MODULE=on

export GOPROXY=https://goproxy.cn

默认端口23456

go build后将文件夹二进制要放在note目录下

nohup ./note &

