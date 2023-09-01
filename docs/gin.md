# Gin

## Gin框架介绍

Go世界里最流行的Web框架，[Github](https://github.com/gin-gonic/gin)上有32K+star。 基于httprouter开发的Web框架。 [中文文档](https://gin-gonic.com/zh-cn/docs/)齐全，简单易用的轻量级框架。

## Gin安装与使用

### 安装

命令行下载就行：

```shell
go get -u github.com/gin-gonic/gin
```

### RESTful API

REST是Representational State Transfer的简称，中文翻译为“表征状态转移”或“表现层状态转化”。

REST的含义就是客户端与Web服务器之间进行交互的时候，使用HTTP协议中的4个请求方法代表不同的动作。

- GET用来获取资源
- POST用来新建资源
- PUT用来更新资源
- DELETE用来删除资源。

只要API程序遵循了REST风格，那就可以称其为RESTful API。目前在前后端分离的架构中，前后端基本都是通过RESTful API来进行交互。

在开发RESTful API的时候通常使用[Postman](https://www.getpostman.com/)来作为客户端的测试工具，大概就是模拟各种请求，以及方便发请求收响应。

## Gin渲染

Gin支持：

- HTML渲染
- 自定义模板渲染
- 静态文件处理
- 使用模板继承
- 补充文件路径处理
- JSON渲染
- XML渲染
- YAML渲染
- protobuf渲染

## Gin获取参数

参数包括：

- 获取querystring参数
- 获取form参数
- 获取JSON参数
- 获取path参数
- 参数绑定

## Gin文件上传

支持：

- 单文件上传
- 多文件上传

注意文件格式以及文件大小。

## Gin重定向

- HTTP重定向
- 路由重定向

## Gin路由与路由组

Web服务可能会对不同路由有不同响应，一个好的方式就是路由分组。Gin支持路由，也支持路由组。

Gin框架中的路由使用的是[httprouter](https://github.com/julienschmidt/httprouter)这个库。

其基本原理就是构造一个路由地址的前缀树。

## Gin中间件

Gin框架允许开发者在处理请求的过程中，加入用户自己的钩子（Hook）函数。这个钩子函数就叫中间件，中间件适合处理一些公共的业务逻辑，比如登录认证、权限校验、数据分页、记录日志、耗时统计等。

### 定义中间件

Gin中的中间件必须是一个gin.HandlerFunc类型。

### 注册中间件

可以对路由，或者路由组，添加一个或者多个中间件。

### 中间件注意事项

- gin默认中间件

gin.Default()默认使用了Logger和Recovery中间件，其中：

Logger中间件将日志写入gin.DefaultWriter，即使配置了GIN_MODE=release。
Recovery中间件会recover任何panic。如果有panic的话，会写入500响应码。
如果不想使用上面两个默认的中间件，可以使用gin.New()新建一个没有任何默认中间件的路由。

- gin中间件中使用goroutine

当在中间件或handler中启动新的goroutine时，不能使用原始的上下文（c *gin.Context），必须使用其只读副本（c.Copy()）。

## Ref

<http://www.ruanyifeng.com/blog/2011/09/restful.html>
<https://www.liwenzhou.com/posts/Go/gin/>
