package router

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/zhj0811/fabric-normal/apiserver/handler"
)

// Router 全局路由
var router *gin.Engine
var onceCreateRouter sync.Once

//GetRouter 获取路由
func GetRouter() *gin.Engine {
	onceCreateRouter.Do(func() {
		router = createRouter()
	})

	return router
}

func createRouter() *gin.Engine {
	r := gin.Default()

	//r.POST("/login", handler.Login)       //登录
	//r.POST("/register", handler.Register) //机构注册
	//r.Use(handler.TokenAuthMiddleware())
	v1 := r.Group("/v1")
	{
		v1.POST("/invoke", handler.Invoke)
		v1.GET("/query", handler.Query)

	}
	return r
}
