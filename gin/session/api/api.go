package api

import (
	"gin/session/middleware"
	"gin/session/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Server(router *gin.Engine) *gin.Engine {
	// 创建用户
	router.POST("/account", CreateAccount)
	// 删除用户
	router.DELETE("/accont", DeleteAccount)
	// 用户登录
	router.GET("/login", Login)
	// 用户退出

	// 业务逻辑
	group := router.Group("/test", middleware.Auth)
	{
		group.GET("/hello", HelloWorld)
	}
	return router
}

// 创建账户
// curl -X POST -d '{"work_id":"123", "password":"123456"}' http://127.0.0.1:8080/account
func CreateAccount(ctx *gin.Context) {
	var account model.Account
	if err := ctx.ShouldBindJSON(&account); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    1001,
			"message": err.Error(),
		})
		return
	}
	// 插入数据库session
	if err := model.CreateAccount(ctx.Request.Context(), &account); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    1002,
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": ""})
}

// 删除账户
func DeleteAccount(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": ""})
}

// 用户登录
// curl -X GET -d '{"work_id":"1212", "password":"123456"}' http://127.0.0.1:8080/login
func Login(ctx *gin.Context) {
	// 绑定参数
	var account model.Account
	if err := ctx.ShouldBindJSON(&account); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    1001,
			"message": err.Error(),
		})
		return
	}

	// 检查用户id对应的session是否存在
	// 存在直接返回
	sess := middleware.GetUserSession(account.WorkId)
	if sess != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"messgae": "",
			"data":    sess,
		})
		return
	}

	// 数据库校验用户密码
	if err := model.LoginAccount(ctx.Request.Context(), &account); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    1001,
			"message": err.Error(),
		})
		return
	}

	// 密码校验通过生成session
	newSess := middleware.NewSession(account.WorkId)
	if newSess == nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    1001,
			"message": "auth service is down",
		})
		return
	}

	// 返回session信息
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "",
		"data":    newSess,
	})
}

// curl -X GET -H "Access-Token:token" http://127.0.0.1:8080/test/hello
func HelloWorld(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "hello world")
}
