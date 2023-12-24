package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

var secrets = []gin.H{
	{"name": "foo", "password": "123456"},
	{"name": "austin", "password": "123456"},
	{"name": "lena", "password": "123456"},
}

var foodData = map[string]Food{}

var ErrorUserNotFound = errors.New("user is not existed")
var ErrorFoodExisted = errors.New("food is existed")
var ErrorFoodNotExisted = errors.New("food is not existed")

type User struct {
	Name     string `form:"name" example:"abc" binding:"required"`
	Password string `form:"password" example:"123456" binding:"required"`
}

type Food struct {
	Name       string `json:"name" example:"拉面" binding:"required"`
	Price      int    `json:"price" example:"123" binding:"required"`
	LeftNum    int    `json:"leftNum" example:"123"`
	UpdateTime string `json:"updateTime" example:"2023-12-11"`
}

func (food *Food) String() string {
	return fmt.Sprintf("name=%s, price=%d, leftNum=%d, updateTime=%s\n", food.Name, food.Price, food.LeftNum, food.UpdateTime)
}

func TestAuthMiddlerWare() {
	router := gin.Default()
	f, _ := os.Create("gin.log")                     // 打开log文件
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout) // 写到log文件
	gin.DefaultErrorWriter = io.MultiWriter(f, os.Stderr)
	adminGroup := router.Group("/admin", Auth()) // 创建路由组设置auth中间件校验
	{
		adminGroup.POST("/food", AddFood)     // 添加food接口
		adminGroup.DELETE("/food", DeletFood) // 删除food接口
	}
	server := &http.Server{Addr: ":8080", Handler: router} // 创建server

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Http server error: %v\n", err)
		}
	}()

	<-signals
	fmt.Println("server receive a signal to shutdown")
	for _, value := range foodData {
		fmt.Println(value.String())
	}
}

func AddFood(ctx *gin.Context) {
	var food Food
	if err := ctx.ShouldBindJSON(&food); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	if _, ok := foodData[food.Name]; ok {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": ErrorFoodExisted.Error()})
		return
	}
	foodData[food.Name] = food
	ctx.String(http.StatusOK, "add record success")
}

func DeletFood(ctx *gin.Context) {
	var food Food
	if err := ctx.ShouldBindJSON(&food); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	if _, ok := foodData[food.Name]; !ok {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": ErrorFoodNotExisted.Error()})
		return
	}
	delete(foodData, food.Name)
	ctx.String(http.StatusOK, "delete record success")
}

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user User
		if err := ctx.ShouldBindQuery(&user); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}
		found := false
		for _, secret := range secrets {
			if secret["name"] == user.Name && secret["password"] == user.Password {
				found = true
				break
			}
		}
		if !found {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": ErrorUserNotFound.Error()})
			return
		}
		ctx.Next()
	}
}
