package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func TestMiddleWare() {
	router := gin.New()
	router.Use(Logger())
	router.GET("/testing", func(ctx *gin.Context) {
		example := ctx.MustGet("example").(string)

		// will print 12345
		log.Println(example)
		ctx.JSON(http.StatusOK, gin.H{"name": "abc", "age": "123"})
	})
	router.Run(":8080")
}

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		t := time.Now()

		// set example variable
		ctx.Set("example", "12345")

		// before request
		ctx.Next() // 调用接口函数

		// after request
		latency := time.Since(t)
		log.Print(latency)

		// access the status we are sending
		status := ctx.Writer.Status()
		log.Println(status)
	}
}
