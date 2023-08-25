<!--
 * @Autor: violet apricity ( Zhuangpx )
 * @Date: 2023-08-23 10:41:34
 * @LastEditors: violet apricity ( Zhuangpx )
 * @LastEditTime: 2023-08-23 11:19:00
 * @FilePath: \docs\服务中心之ETCD.md
 * @Description:  Zhuangpx : Violet && Apricity:/ The warmth of the sun in the winter /
-->
# 服务中心之ETCD

## 前言

我们在做微服务的时候，一个浅显的方式是多个服务有多个service，使用时又会有多个client。
这种简单粗暴的方式显得有些丑陋，那么多service分开，貌似和微服务关系不大了。
有没有一种方式，既能实现服务的独立，又能在一定程度上把它们管理在一起呢？
答案是使用一定的服务中心统一管理，如，本文的ETCD。

> 本文多有基于网络博文，无商业用途，仅为学习

## ETCD介绍

[etcd](https://etcd.io/)

etcd是一个开源的分布式键值对存储工具。在每个coreos节点上面运行的etcd，共同组建了coreos集群的共享数据总线。etcd可以保证coreos集群的稳定，可靠。当集群网络出现动荡，或者当前master节点出现异常时，etcd可以优雅的进行master节点的选举工作，同时恢复集群中损失的数据。

作为使用Go语言开发的一个开源的、高可用的分布式key-value存储系统，ETCD可以用于配置共享和服务的注册和发现。

etcd具有以下特点：

- 完全复制：集群中的每个节点都可以使用完整的存档
- 高可用性：Etcd可用于避免硬件的单点故障或网络问题
- 一致性：每次读取都会返回跨多主机的最新写入
- 简单：包括一个定义良好、面向用户的API（gRPC）
- 安全：实现了带有可选的客户端证书身份验证的自动化TLS
- 快速：每秒10000次写入的基准速度
- 可靠：使用Raft算法实现了强一致、高可用的服务存储目录

## ETCD使用场景

ETCD有很多使用场景，包括：

- 配置管理
- 服务发现
- 选主
- 应用调度
- 分布式队列
- 分布式锁

Etcd的底层数据存储上与一个NoSQL的数据库基本没有差别，但更准确的说法说是一个高可用的键值存储系统。Etcd在设计的初衷主要用于是共享配置和服务发现。

## 对比

### etcd vs redis

先说结论，redis是牺牲数据安全保证快速，而etcd是牺牲速度保证数据安全。

redis 是一种内存数据存储，可用作数据库、缓存或消息代理。redis 支持比 etcd 更广泛的数据类型和结构，并且具有更快的读/写性能。

但是 etcd 具有超强的容错能力、更强的故障转移和持续数据可用性能力，最重要的是，etcd 将所有存储的数据持久化到磁盘，从本质上牺牲了速度以获得更高的可靠性和保证的一致性。由于这些原因，Redis 更适合用作分布式**内存**缓存系统，而不是存储和分布式系统配置信息。

- redis 主从是异步复制的机制，这就导致了其有丢失数据的风险。
- 分布式系统最重要的就是一致性协议，redis 是不支持的。如果发生脑裂，可能两个微服务都会声称自己对某段 IP 的请求负责。

### etcd vs zookeeper

etcd 具有 ZooKeeper 没有的一些重要功能。例如，与 ZooKeeper 不同，etcd 可以执行以下操作：

- 允许动态重新配置集群成员资格。
- 在高负载下执行读/写操作时保持稳定。
- 维护多版本并发控制数据模型。
- 提供可靠的密钥监控，不会在不发出通知的情况下丢弃事件。
- 使用将连接与会话分离的并发原语。
- 支持广泛的语言和框架（ZooKeeper 有自己的自定义 Jute RPC 协议，支持有限的语言绑定）。

Zookeeper有如下缺点

- 复杂。Zookeeper的部署维护复杂，管理员需要掌握一系列的知识和技能；而Paxos强一致性算法也是素来以复杂难懂而闻名于世；另外，Zookeeper的使用也比较复杂，需要安装客户端，官方只提供了java和C两种语言的接口。
- Java编写。这里不是对Java有偏见，而是Java本身就偏向于重型应用，它会引入大量的依赖。而运维人员则普遍希望机器集群尽可能简单，维护起来也不易出错。
- 发展缓慢。Apache基金会项目特有的“Apache Way”在开源界饱受争议，其中一大原因就是由于基金会庞大的结构以及松散的管理导致项目发展缓慢。

etcd的优点：

- 简单。使用Go语言编写部署简单；使用HTTP作为接口使用简单；使用Raft算法保证强一致性让用户易于理解。
- 数据持久化。etcd默认数据一更新就进行持久化。
- 安全。etcd支持SSL客户端安全认证。

### etcd vs consul

Consul 是分布式系统的服务网络解决方案，其功能介于 etcd 和 [Istio 服务网格](https://www.ibm.com/topics/istio) 之间。与 etcd 一样，Consul包含一个基于 Raft 算法的分布式键值存储，并支持 HTTP/JSON 应用程序编程接口（API）。两者都提供动态集群成员配置，但 Consul 对配置数据的多个并发版本没有那么强的控制，并且它可以可靠工作的最大数据库大小更小。

## 在go里使用etcd

以v3为例。

[文档](https://pkg.go.dev/go.etcd.io/etcd/clientv3/concurrency?utm_source=godoc)

> 参考七米老师博文

### 安装

```shell
go get go.etcd.io/etcd/clientv3
```

### put和get操作

put命令用来设置键值对数据，get命令用来根据key获取值。

```go
package main

import (
  "context"
  "fmt"
  "time"

  "go.etcd.io/etcd/clientv3"
)

// etcd client put/get demo
// use etcd/clientv3

func main() {
  //  new point
  cli, err := clientv3.New(clientv3.Config{
    Endpoints:   []string{"127.0.0.1:2379"},
    DialTimeout: 5 * time.Second,
  })
  if err != nil {
    // handle error!
    fmt.Printf("connect to etcd failed, err:%v\n", err)
    return
  }
  fmt.Println("connect to etcd success")
  defer cli.Close()
  // put
  ctx, cancel := context.WithTimeout(context.Background(), time.Second)
  _, err = cli.Put(ctx, "q1mi", "dsb")
  cancel()
  if err != nil {
    fmt.Printf("put to etcd failed, err:%v\n", err)
    return
  }
  // get
  ctx, cancel = context.WithTimeout(context.Background(), time.Second)
  resp, err := cli.Get(ctx, "q1mi")
  cancel()
  if err != nil {
    fmt.Printf("get from etcd failed, err:%v\n", err)
    return
  }
  for _, ev := range resp.Kvs {
    fmt.Printf("%s:%s\n", ev.Key, ev.Value)
  }
}
```

### watch操作

watch用来获取未来更改的通知。

```go

package main

import (
  "context"
  "fmt"
  "time"

  "go.etcd.io/etcd/clientv3"
)

// watch demo

func main() {
  cli, err := clientv3.New(clientv3.Config{
    Endpoints:   []string{"127.0.0.1:2379"},
    DialTimeout: 5 * time.Second,
  })
  if err != nil {
    fmt.Printf("connect to etcd failed, err:%v\n", err)
    return
  }
  fmt.Println("connect to etcd success")
  defer cli.Close()
  // watch key:q1mi change
  rch := cli.Watch(context.Background(), "q1mi") // <-chan WatchResponse
  for wresp := range rch {
    for _, ev := range wresp.Events {
      fmt.Printf("Type: %s Key:%s Value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
    }
  }
}
```

这里是watch一下"qimo"这个key的变化，改变的话就会有通知。

### lease续约

```go
package main

import (
  "fmt"
  "time"
)

// etcd lease

import (
  "context"
  "log"

  "go.etcd.io/etcd/clientv3"
)

func main() {
  cli, err := clientv3.New(clientv3.Config{
    Endpoints:   []string{"127.0.0.1:2379"},
    DialTimeout: time.Second * 5,
  })
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println("connect to etcd success.")
  defer cli.Close()

  // 创建一个5秒的租约
  resp, err := cli.Grant(context.TODO(), 5)
  if err != nil {
    log.Fatal(err)
  }

  // 5秒钟之后, /nazha/ 这个key就会被移除
  _, err = cli.Put(context.TODO(), "/nazha/", "dsb", clientv3.WithLease(resp.ID))
  if err != nil {
    log.Fatal(err)
  }
}
```

### KeepAlive

```go
package main

import (
  "context"
  "fmt"
  "log"
  "time"

  "go.etcd.io/etcd/clientv3"
)

// etcd keepAlive

func main() {
  cli, err := clientv3.New(clientv3.Config{
    Endpoints:   []string{"127.0.0.1:2379"},
    DialTimeout: time.Second * 5,
  })
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println("connect to etcd success.")
  defer cli.Close()

  resp, err := cli.Grant(context.TODO(), 5)
  if err != nil {
    log.Fatal(err)
  }

  _, err = cli.Put(context.TODO(), "/nazha/", "dsb", clientv3.WithLease(resp.ID))
  if err != nil {
    log.Fatal(err)
  }

  // the key 'foo' will be kept forever
  ch, kaerr := cli.KeepAlive(context.TODO(), resp.ID)
  if kaerr != nil {
    log.Fatal(kaerr)
  }
  for {
    ka := <-ch
    fmt.Println("ttl:", ka.TTL)
  }
}
```

## etcd做服务中心进行服务发现

### 什么是服务发现

服务发现其实有两层含义，第一层是实例发现，第二层是端口发现。
如果有两个服务A和B，两个服务都是分布式系统，在某个时刻服务A作为客户端要请求服务B，那么这时从服务A中的某个实例（instance）Ai就要访问服务B中的某个实例Bi，那么Ai找到一个合适的Bi就是服务发现的第一层含义，即实例发现。
Ai找到合适的Bi后其实还没完事，Ai需要知道Bi响应这个请求的应用程序，也就是说需要知道Bi响应这个请求的端口号。由于大部分应用程序的端口号都是随机的，Ai不知道到底请求哪个端口，这就需要Bi找个办法告知Ai。这就是端口发现。

服务发现要解决的也是分布式系统中最常见的问题之一，即在同一个分布式集群中的进程或服务，要如何才能找到对方并建立连接。本质上来说，服务发现就是想要了解集群中是否有进程在监听 udp 或 tcp 端口，并且通过名字就可以查找和连接。

具体来说，服务发现需要实现一下基本功能：

- 服务注册：同一service的所有节点注册到相同目录下，节点启动后将自己的信息注册到所属服务的目录中。
- 健康检查：服务节点定时进行健康检查。注册到服务目录中的信息设置一个较短的TTL，运行正常的服务节点每隔一段时间会去更新信息的TTL ，从而达到健康检查效果。
- 服务发现：通过服务节点能查询到服务提供外部访问的 IP 和端口号。比如网关代理服务时能够及时的发现服务中新增节点、丢弃不可用的服务节点。

### etcd做服务发现

#### 服务注册与监考检查

根据etcd的v3 API，当启动一个服务时候，我们把服务的地址写进etcd，注册服务。同时绑定租约（lease），并以续租约（keep leases alive）的方式检测服务是否正常运行，从而实现健康检查。

主动退出服务时，可以调用Close()方法，撤销租约，从而注销服务。

#### 服务发现

根据etcd的v3 API，很容易想到使用Watch监视某类服务，通过Watch感知服务的添加，修改或删除操作，修改服务列表。

### 结合grpc

结合grpc做微服务的服务注册与负载均衡。

可以参考[这篇文章](https://juejin.cn/post/7232483262979719205)

## Ref

<https://www.liwenzhou.com/posts/Go/etcd/>
<https://www.wyx.cloudns.asia/blog/2021/07/10/service_discovery>
<https://laravelacademy.org/post/21218>
<https://guanhonly.github.io/2020/08/30/etcd-service-discovery/>
<https://juejin.cn/post/7101947466722836487>
