<!--
 * @Autor: violet apricity ( Zhuangpx )
 * @Date: 2023-08-24 22:44:42
 * @LastEditors: violet apricity ( Zhuangpx )
 * @LastEditTime: 2023-08-25 23:04:07
 * @FilePath: \docs\viper_yaml.md
 * @Description:  Zhuangpx : Violet && Apricity:/ The warmth of the sun in the winter /
-->

# viper_yaml

> 配置入门-viper与yaml

我们在做一些服务的时候，例如，我们需要开放一些端口，如我我们写成静态的，在项目开发的某个时期，我们想要有一个端口的改变，或者提供端口自定义的行为，这时代码改动就比较大。正确的方式是写成配置文件，合理配置。

一个友好的方式是，使用viper配置，配置文件，如，可以使用yaml。

## yaml入门

### 简介

YAML 是 "YAML Ain't a Markup Language"（YAML 不是一种标记语言）的递归缩写。在开发的这种语言时，YAML 的意思其实是："Yet Another Markup Language"（仍是一种标记语言）。

YAML 的语法和其他高级语言类似，并且可以简单表达清单、散列表，标量等数据形态。它使用空白符号缩进和大量依赖外观的特色，特别适合用来表达或编辑数据结构、各种配置文件、倾印调试内容、文件大纲（例如：许多电子邮件标题格式和YAML非常接近）。

YAML 的配置文件后缀为 .yml，如：xxx.yml 。

基本语法：

- 大小写敏感
- 使用缩进表示层级关系
- 缩进不允许使用tab，只允许空格
- 缩进的空格数不重要，只要相同层级的元素左对齐即可
- '#'表示注释

数据类型：

- 对象：键值对的集合，又称为映射（mapping）/ 哈希（hashes） / 字典（dictionary）
- 数组：一组按次序排列的值，又称为序列（sequence） / 列表（list）
- 纯量（scalars）：单个的、不可再分的值

### yaml对象

对象的一组键值对，使用冒号结构表示。

```yaml
animal: pets
```

Yaml 也允许另一种写法，将所有键值对写成一个行内对象。

```yaml
hash: { name: Steve, foo: bar }
```

### yaml数组

一组连词线开头的行，构成一个数组。

```yaml
- Cat
- Dog
- Goldfish
```

数据结构的子成员是一个数组，则可以在该项下面缩进一个空格。

```yaml
-
 - Cat
 - Dog
 - Goldfish
```

数组也可以采用行内表示法。

```yaml
animal: [Cat, Dog]
```

### yaml复合结构

就是把对象和数组结合使用，形成复合结构。

### 纯量

类似面量。纯量是最基本的、不可再分的值。包括：

```markdown
- 字符串
- 布尔值
- 整数
- 浮点数
- Null
- 时间
- 日期
```

> PS：强制类型转换
> yaml允许使用两个感叹号，强制转换数据类型。
> `e: !!str 123`

### 字符串

- 字符串默认不使用引号表示。
- 如果字符串之中包含空格或特殊字符，需要放在引号之中。
- 单引号和双引号都可以使用，双引号不会对特殊字符转义。
- 单引号之中如果还有单引号，必须连续使用两个单引号转义。
- 字符串可以写成多行，从第二行开始，必须有一个单空格缩进。换行符会被转为空格。
- 多行字符串可以使用|保留换行符，也可以使用>折叠换行。
- +表示保留文字块末尾的换行，-表示删除字符串末尾的换行。

### 引用

有一个`&`表示锚点，一个`*`表示别名。如：

```yaml
defaults: &defaults
  adapter:  postgres
  host:     localhost

development:
  database: myapp_development
  <<: *defaults

test:
  database: myapp_test
  <<: *defaults
```

等价于：

```yaml
defaults:
  adapter:  postgres
  host:     localhost

development:
  database: myapp_development
  adapter:  postgres
  host:     localhost

test:
  database: myapp_test
  adapter:  postgres
  host:     localhost
```

### 函数和正则表达式的转换

可以把函数和正则表达式转为字符串。

## viper入门

Viper是适用于Go应用程序的完整配置解决方案。它被设计用于在应用程序中工作，并且可以处理所有类型的配置需求和格式。

### 安装

```shell
go get github.com/spf13/viper
```

### viper简介

Viper是适用于Go应用程序（包括Twelve-Factor App）的完整配置解决方案。它被设计用于在应用程序中工作，并且可以处理所有类型的配置需求和格式。它支持以下特性：

- 设置默认值
- 从JSON、TOML、YAML、HCL、envfile和Java properties格式的配置文件读取配置信息
- 实时监控和重新读取配置文件（可选）
- 从环境变量中读取
- 从远程配置系统（etcd或Consul）读取并监控配置变化
- 从命令行参数读取配置
- 从buffer读取配置
- 显式配置值

### api

viper的api，官网写的很清楚，不赘述。

### 例子

比如有个配置文件`config.yaml`。

#### 直接使用viper管理配置

在gin框架搭建的web项目中使用viper，使用viper加载配置文件中的信息，并在代码中直接使用viper.GetXXX()方法获取对应的配置值。

```go
package main

import (
  "fmt"
  "net/http"

  "github.com/gin-gonic/gin"
  "github.com/spf13/viper"
)

func main() {
  viper.SetConfigFile("./conf/config.yaml") // 指定配置文件路径
  err := viper.ReadInConfig()        // 读取配置信息
  if err != nil {                    // 读取配置信息失败
    panic(fmt.Errorf("Fatal error config file: %s \n", err))
  }

  // 监控配置文件变化
  viper.WatchConfig()

  r := gin.Default()
  // 访问/version的返回值会随配置文件的变化而变化
  r.GET("/version", func(c *gin.Context) {
    c.String(http.StatusOK, viper.GetString("version"))
  })

  if err := r.Run(
    fmt.Sprintf(":%d", viper.GetInt("port"))); err != nil {
    panic(err)
  }
}
```

#### 使用结构体变量保存配置信息

可以在项目中定义与配置文件对应的结构体，viper加载完配置信息后使用结构体变量保存配置信息。

```go
package main

import (
  "fmt"
  "net/http"

  "github.com/fsnotify/fsnotify"

  "github.com/gin-gonic/gin"
  "github.com/spf13/viper"
)

type Config struct {
  Port    int    `mapstructure:"port"`
  Version string `mapstructure:"version"`
}

var Conf = new(Config)

func main() {
  viper.SetConfigFile("./conf/config.yaml") // 指定配置文件路径
  err := viper.ReadInConfig()               // 读取配置信息
  if err != nil {                           // 读取配置信息失败
    panic(fmt.Errorf("Fatal error config file: %s \n", err))
  }
  // 将读取的配置信息保存至全局变量Conf
  if err := viper.Unmarshal(Conf); err != nil {
    panic(fmt.Errorf("unmarshal conf failed, err:%s \n", err))
  }
  // 监控配置文件变化
  viper.WatchConfig()
  // 注意！！！配置文件发生变化后要同步到全局变量Conf
  viper.OnConfigChange(func(in fsnotify.Event) {
    fmt.Println("夭寿啦~配置文件被人修改啦...")
    if err := viper.Unmarshal(Conf); err != nil {
      panic(fmt.Errorf("unmarshal conf failed, err:%s \n", err))
    }
  })

  r := gin.Default()
  // 访问/version的返回值会随配置文件的变化而变化
  r.GET("/version", func(c *gin.Context) {
    c.String(http.StatusOK, Conf.Version)
  })

  if err := r.Run(fmt.Sprintf(":%d", Conf.Port)); err != nil {
    panic(err)
  }
}
```

## Ref

<https://www.ruanyifeng.com/blog/2016/07/yaml.html>
<https://www.runoob.com/w3cnote/yaml-intro.html>
<https://www.liwenzhou.com/posts/Go/viper/>
