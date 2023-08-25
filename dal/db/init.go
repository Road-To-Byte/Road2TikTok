/*
 * @Autor: violet apricity ( Zhuangpx )
 * @Date: 2023-08-23 19:52:24
 * @LastEditors: violet apricity ( Zhuangpx )
 * @LastEditTime: 2023-08-25 16:12:07
 * @FilePath: \Road2TikTok\dal\db\init.go
 * @Description:  Zhuangpx : Violet && Apricity:/ The warmth of the sun in the winter /
 */
package db

import (
	"fmt"
	"log"
	"time"

	"github.com/Road-To-Byte/Road2TikTok/pkg/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

//	db的初始化操作

// 全局变量
var (
	__err  error
	__db   *gorm.DB
	config = viper.Init("db")
)

// 获取gorm.DB对象
func GetDB() *gorm.DB {
	return __db
}

// DSN:Data Source Name
// eg: <driver-name>://<username>:<password>@<host>:<port>/<database-name>
// 根据 driverWithRole 获取 DSN
func getDSN(driverWithRole string) string {
	username := config.Viper.GetString(fmt.Sprintf("%s.username", driverWithRole))
	password := config.Viper.GetString(fmt.Sprintf("%s.password", driverWithRole))
	host := config.Viper.GetString(fmt.Sprintf("%s.host", driverWithRole))
	port := config.Viper.GetInt(fmt.Sprintf("%s.port", driverWithRole))
	DBname := config.Viper.GetString(fmt.Sprintf("%s.database", driverWithRole))
	//	charset=utf8
	//	parseTime=True
	//	loc=Local	本地时区
	arg := "?" + "charset=utf8" + "&" + "parseTime=True" + "&" + "loc=Local"
	//	dsn
	//                  u  p      h  p   D a
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s%s", username, password, host, port, DBname, arg)
	return dsn
}

// 初始化
// 这里用到了gorm的一个插件dbresolver
// 用于实现数据库读写分离和负载均衡功能
// 这个插件允许配置多个数据库连接（主数据库和多个从数据库），并且可以根据一定的策略将读操作分发到不同的从数据库上，实现读操作的负载均衡
func init() {
	//	主数据库 sources
	dsn1 := getDSN("mysql.source")
	//	通过gorm和mysql创建连接
	__db, __err := gorm.Open(mysql.Open(dsn1), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Info),
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	})
	if __err != nil {
		panic(__err.Error())
	}
	//	从数据库 replicas
	dsn2 := getDSN("mysql.replica1")
	dsn3 := getDSN("mysql.replica2")
	// 配置 dbresolver 注意区分主数据库和从数据库
	__db.Use(dbresolver.Register(dbresolver.Config{
		Sources:  []gorm.Dialector{mysql.Open(dsn1)},
		Replicas: []gorm.Dialector{mysql.Open(dsn2), mysql.Open(dsn3)},
		//	负载均衡策略： 随机(简单粗暴)
		Policy: dbresolver.RandomPolicy{},
		// print sources/replicas mode in logger
		TraceResolverMode: false,
	}))
	//	看见auto 没错 这是和自动化有关的
	//	AutoMigrate会自动检查已定义的模型（这里是结构体），然后在数据库里创建，以及更新相应的表结构，从而匹配这些模型
	//	所以他能补一些缺失的表，外键，约束，列和索引等 不过不用担心 他不会删除（谁给你这么大权利.jpg）
	// err := __db.AutoMigrate(&User{}, &Video{}, &Comment{}, &FavoriteVideoRelation{}, &FollowRelation{}, &Message{}, &FavoriteCommentRelation{})
	var err error
	if err != nil {
		log.Fatalln(err.Error())
	}
	//	拿到sql对象
	db, err := __db.DB()
	if err != nil {
		log.Fatalln(err.Error())
	}
	//	最大使用连接数
	db.SetMaxOpenConns(1000)
	//	最大闲置连接数
	db.SetMaxIdleConns(30)
	//	连接的最长声明周期
	db.SetConnMaxLifetime(60 * time.Minute)
}
