package main

import (
	"context"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ServerHttp(ctx context.Context, address string, cal *Calculator) error {
	lisAddr, err := net.Listen("tcp", address)
	if nil != err {
		return err
	}
	r := gin.Default()
	r.GET("/Add", func(c *gin.Context) {
		var task AddTask
		if err := c.ShouldBind(&task); err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		task.result = make(chan int, 1)
		cal.Submit(&task)
		result := task.Wait()
		c.JSON(http.StatusOK, gin.H{
			"code":   http.StatusOK,
			"result": result,
		})
	})
	r.GET("/Sub", func(c *gin.Context) {
		var task SubTask
		if err := c.ShouldBind(&task); err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		task.result = make(chan int, 1)
		cal.Submit(&task)
		result := task.Wait()
		c.JSON(http.StatusOK, gin.H{
			"code":   http.StatusOK,
			"result": result,
		})
	})
	go func() {
		<-ctx.Done()
		lisAddr.Close()
	}()
	return r.RunListener(lisAddr)
}
