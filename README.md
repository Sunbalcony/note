# 一个web笔记本

本站demo：长期维护 https://note.grpc.fun

## 更新记录

2021/09/10

添加通过API POST请求创建

2021/09/30

配置项解耦，全部放在config/application.yml中

2022/03/25

增加对接redis存储 配置文件type项目 项目工程化处理 优化tag生成逻辑

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

export GO111MODULE=on

export GOPROXY=https://goproxy.cn

go build 保持 conf static note二进制文件同步录下

修改配置文件

nohup ./note &

### docker

本地创建配置文件路径 mkdir -p /data/note/conf

修改配置文件并保存
```shell
cat > /data/note/conf/application.yml << EOF
note:
  serverPort: 8080
  keylength: 6
  type: 1
  mysqlUrl: mysql.com:3306
  mysqlUsername: root
  mysqlPassword: 123456
  mysqlDatabasename: notes
  timezone: Asia/Shanghai
  redisUrl: 192.168.8.8:6379
  redisPassword: xayf
  redisDatabaseNum: 1
EOF
```

docker run -dit -v /data/note/conf:/root/conf --name note -p 8080:8080 sooemma/note:1.0


