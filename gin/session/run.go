package session

import (
	"context"
	"gin/session/api"
	"gin/session/dbstore"
	"gin/session/rds"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func TestSession(listenAddr string) {
	// 初始化redis连接
	rds.InitCache("10.177.54.121:6379", "123456")
	// 初始化mysql连接
	dbstore.InitDBStore()

	router := gin.Default()
	// 日志写入文件
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	gin.DefaultErrorWriter = io.MultiWriter(f, os.Stderr)

	// 优雅退出
	server := &http.Server{Addr: listenAddr, Handler: api.Server(router)}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Http server error: %v\n", err)
		}
	}()

	<-signals
	log.Println("shutting down server...")

	ctx, cfn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cfn()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}
	log.Println("Server exiting")
}
