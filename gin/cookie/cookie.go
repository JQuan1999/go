package cookie

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func TestCookie() {
	router := gin.Default()

	router.GET("/cookie", func(ctx *gin.Context) {
		cookie, err := ctx.Cookie("gin_cookie")
		if err != nil {
			cookie = "Not set"
			ctx.SetCookie("gin_cookie", "test", 3600, "/", "localhost", false, true)
		}
		fmt.Printf("Cookie value: %s \n", cookie)
	})

	router.Run(":8080")
}
