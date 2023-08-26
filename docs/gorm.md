<!--
 * @Autor: violet apricity ( Zhuangpx )
 * @Date: 2023-08-26 22:07:29
 * @LastEditors: violet apricity ( Zhuangpx )
 * @LastEditTime: 2023-08-26 22:19:21
 * @FilePath: \docs\gorm.md
 * @Description:  Zhuangpx : Violet && Apricity:/ The warmth of the sun in the winter /
-->
# gorm

我们在写go的时候，难免需要操作数据库，关于操作数据库有一种很友好的方式，叫ORM，Object Relational Mapping(对象关系映射)，也就是把面向对象的概念和数据库中表的概念对应起来。

在go里也有这种方式，一个常见的框架，即本文的gorm。

## gorm介绍

gorm是一个使用Go语言编写的ORM框架。它文档齐全，对开发者友好，支持主流数据库。

[官网](https://github.com/go-gorm/gorm)

[中文文档](https://gorm.io/zh_CN/)

## 安装

以mysql为例。

```shell
//安装MySQL驱动
go get -u gorm.io/driver/mysql
//安装gorm包
go get -u gorm.io/gorm
```

## 连接数据库

需要先导入驱动程序，go非常友好地支持了这一操作，直接import即可。

## api

官网文档介绍，不再赘述。

## GORM Model

在使用ORM工具时，通常我们需要在代码中定义模型（Models）与数据库中的数据表进行映射，在GORM中模型（Models）通常是正常定义的结构体、基本的go类型或它们的指针。 同时也支持sql.Scanner及driver.Valuer接口（interfaces）。

为了方便模型定义，GORM内置了一个gorm.Model结构体。gorm.Model是一个包含了ID, CreatedAt, UpdatedAt, DeletedAt四个字段的Golang结构体。

```go
// gorm.Model 定义
type Model struct {
  ID        uint `gorm:"primary_key"`
  CreatedAt time.Time
  UpdatedAt time.Time
  DeletedAt *time.Time
}
```

可以将它嵌入到你自己的模型中：

```go
// 将 `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`字段注入到`User`模型中
type User struct {
  gorm.Model
  Name string
}
```

也可以完全自己定义模型：

```go
// 不使用gorm.Model，自行定义模型
type User struct {
  ID   int
  Name string
}
```

一个例子：

```go
type User struct {
  gorm.Model
  Name         string
  Age          sql.NullInt64
  Birthday     *time.Time
  Email        string  `gorm:"type:varchar(100);unique_index"`
  Role         string  `gorm:"size:255"` // 设置字段大小为255
  MemberNumber *string `gorm:"unique;not null"` // 设置会员号（member number）唯一并且不为空
  Num          int     `gorm:"AUTO_INCREMENT"` // 设置 num 为自增类型
  Address      string  `gorm:"index:addr"` // 给address字段创建名为addr的索引
  IgnoreMe     int     `gorm:"-"` // 忽略本字段
}
```

> 使用结构体声明模型时，标记（tags）是可选项。

## 主键、表名、列名的约定

### 主键（Primary Key）

GORM 默认会使用名为ID的字段作为表的主键。

### 表名（Table Name）

表名默认就是结构体名称的复数，也可以通过Table()指定表名。

### 列名（Column Name）

列名由字段名称进行下划线分割来生成：

```go
type User struct {
  ID        uint      // column name is `id`
  Name      string    // column name is `name`
  Birthday  time.Time // column name is `birthday`
  CreatedAt time.Time // column name is `created_at`
}
```

可以使用结构体tag指定列名：

```go
type Animal struct {
  AnimalId    int64     `gorm:"column:beast_id"`         // set column name to `beast_id`
  Birthday    time.Time `gorm:"column:day_of_the_beast"` // set column name to `day_of_the_beast`
  Age         int64     `gorm:"column:age_of_the_beast"` // set column name to `age_of_the_beast`
}
```

## 时间戳跟踪

### CreatedAt

如果模型有 CreatedAt字段，该字段的值将会是初次创建记录的时间。

db.Create(&user) // `CreatedAt`将会是当前时间

```go
// 可以使用`Update`方法来改变`CreateAt`的值
db.Model(&user).Update("CreatedAt", time.Now())
```

### UpdatedAt

如果模型有UpdatedAt字段，该字段的值将会是每次更新记录的时间。

```go
db.Save(&user) // `UpdatedAt`将会是当前时间

db.Model(&user).Update("name", "jinzhu") // `UpdatedAt`将会是当前时间
```

### DeletedAt

如果模型有DeletedAt字段，调用Delete删除该记录时，将会设置DeletedAt字段为当前时间，而不是直接将记录从数据库中删除。

## Ref

<https://mszlu.com/go/gorm/01/01.html#_1-%E5%85%A5%E9%97%A8>
<https://www.liwenzhou.com/posts/Go/gorm/>
