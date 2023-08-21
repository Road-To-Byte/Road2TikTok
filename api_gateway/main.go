/*
 * @Autor: violet apricity ( Zhuangpx )
 * @Date: 2023-08-14 18:04:38
 * @LastEditors: violet apricity ( Zhuangpx )
 * @LastEditTime: 2023-08-21 01:29:16
 * @FilePath: \Road2TikTok\api_gateway\main.go
 * @Description:  Zhuangpx : Violet && Apricity:/ The warmth of the sun in the winter /
 */
package main

import (
	"github.com/Road-To-Byte/Road2TikTok/api_gateway/handler"
	"github.com/gin-gonic/gin"
)

func InitGin() *gin.Engine {
	gn := gin.Default()
	return gn
}

//	路由分组 层级api
func registerGroup(r *gin.Engine) {
	douyin := r.Group("/douyin")
	{
		//	user api
		user := douyin.Group("/user")
		{
			//	用户信息
			user.GET("/", handler.UserInfo)
			//	用户注册
			user.POST("/register/")
			//	用户登录
			user.POST("/login/")
		}
		//	message api
		message := douyin.Group("/message")
		{
			//	聊天记录
			message.GET("/chat/")
			//	消息操作
			message.POST("/action/")
		}
		//	relation api
		relation := douyin.Group("/relation")
		{
			//	关系操作
			relation.POST("/action/")
			//	用户好友列表
			relation.GET("/friend/list/")
			//	用户关注列表
			relation.GET("follow/list/")
			//	用户粉丝列表
			relation.GET("/follower/list/")
		}
		//	publish api
		publish := douyin.Group("/publish/")
		{
			//	发布列表
			publish.GET("/list/")
			//	视频投稿
			publish.POST("/action/")
		}
		//	feed api
		feed := douyin.Group("/feed")
		{
			feed.GET("/")
		}
		//	favorite api
		favorite := douyin.Group("/favorite/")
		{
			//	喜欢列表
			favorite.GET("/list/")
			//	赞操作
			favorite.POST("/action/")
		}
		//	comment api
		comment := douyin.Group("/comment")
		{
			//	视频评论列表
			comment.GET("/list/")
			//	评论操作
			comment.POST("/action/")
		}
	}
}

func main() {
	//	初始化
	router := InitGin()
	//	路由分组 api注册
	registerGroup(router)
	//	启动
	router.Run(":9090")
}
