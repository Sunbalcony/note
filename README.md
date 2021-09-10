# 一个web笔记本

本站demo：长期维护 https://note.valarx.com

## 更新记录

2021/09/10

添加通过API POST请求创建

```
 创建：
 POST /api/create 
 参数 Text 
 通过API来创建一个笔记本
 返回 一个URL 打开即可访问
 
 更新
 POST /api/update
 参数 Tag Text
 通过API来更新一个已经存在笔记本的内容
 
```

## Nginx代理配置

```shell

        location / {
                proxy_pass http://127.0.0.1:23456;
                #携带域名
                proxy_set_header       Host $host;
                
        }

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

