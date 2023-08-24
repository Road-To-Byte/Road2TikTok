<!--
 * @Autor: violet apricity ( Zhuangpx )
 * @Date: 2023-08-10 20:52:30
 * @LastEditors: violet apricity ( Zhuangpx )
 * @LastEditTime: 2023-08-21 01:01:55
 * @FilePath: \Road2TikTok\zpx.md
 * @Description:  Zhuangpx : Violet && Apricity:/ The warmth of the sun in the winter
-->
# zpx

## Pre list

- Gin web框架
- Hertz web框架，字节跳动自己研发
- Kitex rpc框架
- etcd 分布式键值存储，服务注册与发现
- Gorm 操作数据库的对象关系映射(ORM)库
- Redis 数据存储
- JWT token的生成与校验
- 消息队列

HTTP框架用gin
rpc/微服务框架用grpc 并用etcd做服务注册和发现
dal数据层

## cmd

### proto

```shell
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pb/xx.proto
```

## package

### pkg

- jwt

  - ```shell
    go get go.uber.org/zap
    ```

- viper

  - ```shell
    go get github.com/spf13/viper
    ```

- zap

  - ```shell
    go get go.uber.org/zap
    ```

  - lumberjack

    - ```shell
      go get gopkg.in/natefinch/lumberjack.v2
      ```

- etcd

  - ```shell
    go get go.etcd.io/etcd/clientv3
    或者
    go get go.etcd.io/etcd/client/v3@v3.5.9
    ```

- x

### dal

- gorm

  - ```shell
    go get -u gorm.io/gorm
    ```
    
  - plugin
  
    - ```shell
      go get gorm.io/plugin/dbresolver
      ```
  
- mysql

  - ```shell
    go get -u gorm.io/driver/mysql
    ```

- 

## Version

- mysql：8.0.29
- go：1.18.4

## Ref

- <https://github.com/bytedance-youthcamp-jbzx/tiktok/>
- <https://github.com/Go-To-Byte/DouSheng/>
- <https://github.com/writiger/dousheng_server/>
- <https://github.com/a76yyyy/tiktok/>
- [一个基于 gin+ grpc + etcd 等框架开发的小栗子](https://www.cnblogs.com/M-Anonymous/p/17159371.html)
