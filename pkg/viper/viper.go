/*
 * @Autor: violet apricity ( Zhuangpx )
 * @Date: 2023-08-22 16:29:00
 * @LastEditors: violet apricity ( Zhuangpx )
 * @LastEditTime: 2023-08-23 20:08:33
 * @FilePath: \Road2TikTok\pkg\viper\viper.go
 * @Description:  Zhuangpx : Violet && Apricity:/ The warmth of the sun in the winter /
 */
package viper

import (
	"log"

	"github.com/spf13/viper"
)

//	Config 结构
type Config struct {
	Viper *viper.Viper
}

//	Init 根据配置文件的文件名 生成一个新的Config
func Init(configName string) Config {
	//	new 一个 Config
	config := Config{Viper: viper.New()}
	v := config.Viper
	//	设置文件类型
	v.SetConfigType("yml")
	//	设置文件名
	v.SetConfigName(configName)
	//	设置搜索路径 可以添加多条
	v.AddConfigPath("./config")
	v.AddConfigPath("../config")
	v.AddConfigPath("../../config")
	//	读取配置文件
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Viper config errno: %v", err)
	}
	return config
}
