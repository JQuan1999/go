

# gin

### curl

- `-X <method>`：指定HTTP请求的方法，如GET、POST、PUT、DELETE等。
- `-H <header>`：指定请求头，如`-H "Content-Type: application/json"`。
- `-d <data>`：指定请求体中要发送的数据，如`-d '{"name": "Alice"}'`。
- `-u <user:password>`：指定HTTP认证的用户名和密码，如`-u "admin:123456"`。
- `-v`：显示请求和响应的详细信息，包括请求头、响应头、响应体等。通过这个选项可以用来调试Web应用程序和服务端。
- `-o <file>`：将响应内容保存到指定的文件中，如`-o "response.json"`。

```
curl https://www.example.com // 发送get请求
curl -X POST -H "Content-Type: application/json" -d '{"name": "Alice"}' https://www.example.com // 发送post请求
curl -u "admin:123456" https://www.example.com // 发送http认证请求
curl -o response.json https://www.example.com // 保存响应体到文件中
```

### resetful API

使用GET, POST, PUT, PATCH, DELETE and OPTIONS

`````go
func main() {
    router := gin.Default()
    router.GET("/someGet", getting)
    router.POST("/somePost", posting)
    router.PUT("/somePut", putting)
 	router.DELETE("/someDelete", deleting)
  	router.PATCH("/somePatch", patching)
  	router.HEAD("/someHead", head)
  	router.OPTIONS("/someOptions", options)
	
    router.Run()
}
`````

### 路径参数

````go
func main() {
    router := gin.Default()
    
    // handler match /user/john will not match /user or /user
    router.GET("/user/:name", func(c *gin.Context){
        name := c.Param("name")
        action := c.Param("action")
        message := name + " is " + action
        c.String(http.StatusOk, message)
    })
    
    // for each method request context will hold the route definetion
    router.POST("/user/:name/*action", func(c *gin.Context){
        b := c.FullPath() == "/user/:name/*action" // true
        c.String(http.StatusOK, "%t", b)
    })
    
    // for each method request context will hold the route definition
    router.GET("/user/groups", func(c *gin.Context) {
    	c.String(http.StatusOK, "The available groups are [...]")
  	})
    
    router.Run(":8080")
}
````

### 获取query参数

````go
func main() {
	router := gin.Default()
	
	// 解析query string
	router.GET("/welcome", func(c *gin.Context){
		firstName := c.DefaultQuery("firstName", "Guest")
		lastName := c.Query("lastName")
		c.String(http.StatusOk, "Hello %s %s", firstName, lastName)
	})
	router.Run(":8080")
}
````

### 获取请求中的表单数据

````go
func main() {
	router := gin.Default()
	router.POST("/form_post", func(c *gin.Context){
		message := c.PostForm("message") // PostForm()方法来获取POST请求中的表单数据，该方法接受一个参数，即表单中的字段名
		nick := c.DefaultPostForm("nick", "anonymous") // 用于获取POST请求中的表单数据，如果指定的表单字段不存在，则返回默认值anonymous
		
		c.JSON(http.StatusOK, gin.H{
			"status": "posted",
			"message": "message",
			"nick": nick,
		})
	})
	router.Run(":8080")
}
````

### Query+Post Form

````go
POST /post?id=1234&page=1 HTTP/1.1
Content-Type: application/x-www-form-urlencoded

name=manu&message=this_is_great
````

````go
func main() {
  router := gin.Default()

  router.POST("/post", func(c *gin.Context) {

    id := c.Query("id")
    page := c.DefaultQuery("page", "0")
    name := c.PostForm("name")
    message := c.PostForm("message")

    fmt.Printf("id: %s; page: %s; name: %s; message: %s", id, page, name, message)
    // id: 1234; page: 1; name: manu; message: this_is_great
  })
  router.Run(":8080")
}
````

### 上传文件

````go
func main() {
	router := gin.Default()
	// set lower memory
	router.MaxMultipartMemory = 8 << 20
	router.Post("/upload", func(c *gin.Context){
		// single file
		file, _ := c.FormFile("file")
		log.Println(file.Filename)
		
		// upload the file to specific dst
		c.SaveUploadedFile(file, dst)
		
		c.String(http.StatusOk, fmt.Sprintf("'%s' uploaded!", file.FileName))
	})
	router.Run(":8080")
}
````

````
// 使用curl
curl -X POST http://localhost:8080/upload \
  -F "file=@/Users/appleboy/test.zip" \
  -H "Content-Type: multipart/form-data"
````

### 上传多个文件

````go
func main() {
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20
	router.POST("/upload", func(c *gin.Context){
		form, _ := c.MultipartForm()
		files := form.File["upload[]"]
		for _, file := range files {
			log.Println(file.Filename)
			
			// upload the file to specific dst.
			c.SaveUploadFile(file, dst)
		}
		c.String(http.StatusOk, fmt.Sprintf("%d files uploaded!", len(files)))
	})
	router.Run(":8080")
}
// curl
curl -X POST http://localhost:8080/upload \
  -F "upload[]=@/Users/appleboy/test1.zip" \
  -F "upload[]=@/Users/appleboy/test2.zip" \
  -H "Content-Type: multipart/form-data"
````

### 路由组

```go
func main() {
	router := gin.Default()
	// simple group v1
	v1 := router.Group("/v1")
	{
		v1.POST("/login", loginEndpoint)
		v1.POST("/submit", submitEndpoint)
		v1.POST("/read", readEndpoint)
	}
	// simple group v2
	v2 := router.Group()
	{
		v2.POST("/login", loginEndpoint)
		v2.POST("/submit", submitEndpoint)
    	v2.POST("/read", readEndpoint)
	}
	router.Run(":8080")
}
```

### 空白gin默认不使用中间件

```go
r := gin.New()
r := gin.Default() // 使用了默认附带的logger和recovery中间件
```

### 使用中间件

```go
func main() {
	// Creates a router without any middleware by default
  	r := gin.New()
  	// Global middleware
  	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
  	// By default gin.DefaultWriter = os.Stdout
 	r.Use(gin.Logger())
 	
 	// Recovery middleware recovers from any panics and writes a 500 if there was one.
 	r.Use(gin.Recovery())
 	
 	// Per route middleware, you can add as many as you desire.
  	r.GET("/benchmark", MyBenchLogger(), benchEndpoint)
  	
  	// Authorization group
  	// authorized := r.Group("/", AuthRequired())
  	// exactly the same as:
  	authorized := r.Group("/")
  	
  	// per group middleware! in this case we use the custom created
 	// AuthRequired() 中间价只在authorized组中生效
	authorized.Use(AuthRequired())
	{
		authorized.POST("/login", loginEndpoint)
		authorized.POST("/submit", submitEndpoint)
    	authorized.POST("/read", readEndpoint)
    	
    	// 嵌套group
    	testing := authorized.Group("testing")
    	// visit 0.0.0.0:8080/testing/analytics
    	testing.GET("/analytics", analyticsEndpoint)
	}
	// Listen and serve on 0.0.0.0:8080
 	r.Run(":8080")
}
```

### 自定义恢复函数

```go
func main() {
	// Creates a router without any middleware by default
  	r := gin.New()
  	
  	r.Use(gin.Logger())
  	// Recovery middleware recovers from any panics and writes a 500 if there was one.
  	r.Use(gin.CustomRecovery(func(c *gin.Context, recovered any) {
        if err, ok := recovered.(string); ok {
          c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
        }
    	c.AbortWithStatus(http.StatusInternalServerError)
  	}))
  	
  	r.GET("/panic", func(c *gin.Context){
  		// panic with a string -- the custom middleware could save this to a database or report it to the user
  		panic("foo")
  	})
  	
  	r.GET("/", func(c *gin.Context) {
    	c.String(http.StatusOK, "ohai")
  	})

  // Listen and serve on 0.0.0.0:8080
  r.Run(":8080")
}
```

### 自定义log文件

```go
func main() {
  // Disable Console Color, you don't need console color when writing the logs to file.
  gin.DisableConsoleColor()

  // Logging to a file.
  f, _ := os.Create("gin.log")
  gin.DefaultWriter = io.MultiWriter(f)

  // Use the following code if you need to write the logs to file and console at the same time.
  // gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

  router := gin.Default()
  router.GET("/ping", func(c *gin.Context) {
      c.String(http.StatusOK, "pong")
  })

  router.Run(":8080")
}
```

### 模型绑定和验证

可以绑定json、xml、yaml、toml和标准的表单值(foo=bar&&boo=baz)

* Must Bind
  * 方法：`Bind`,`BindJson`,`BindXml`, `BindQuery`,`BindHeader`,`BindToml`
  * 行为：`MustBindWith`绑定失败时，这个请求将会abort：`c.AbortWithError(400, err).SetType(ErrorTypeBind)`. 该方法将会终止请求处理，并返回一个错误响应，因此在使用该方法之后，不能再对响应进行任何操作。
  * 
* ShouldBind
  * ``ShouldBind`, `ShouldBindJSON`, `ShouldBindXML`, `ShouldBindQuery`, `ShouldBindYAML`, `ShouldBindHeader`, `ShouldBindTOML`,`ShouldBindQuery` 只绑定query参数不绑定post数据。

可以指定`binding:"required"`绑定必须的字段，当绑定时出现空值会报错。

````go
// Binding from JSON
type Login struct {
	User     string `form:"user" json:"user" xml:"user"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

func TestShouldBind() {
	router := gin.Default()

	// binding json
	router.POST("/loginJSON", func(c *gin.Context) {
		var json Login
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if json.User != "manu" || json.Password != "123" {
			c.JSON(http.StatusOK, gin.H{"status": "unauthorized"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
	})

	// Example for binding XML (
	//  <?xml version="1.0" encoding="UTF-8"?>
	//  <root>
	//    <user>manu</user>
	//    <password>123</password>
	//  </root>)
	router.POST("/loginXML", func(c *gin.Context) {
		var xml Login
		if err := c.ShouldBindJSON(&xml); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if xml.User != "manu" || xml.Password != "123" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
	})

	// Example for binding a HTML form (user=manu&password=123)
	router.POST("/loginForm", func(c *gin.Context) {
		var form Login
		// This will infer what binder to use depending on the content-type header.
		if err := c.ShouldBind(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if form.User != "manu" || form.Password != "123" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
	})

	// Listen and serve on 0.0.0.0:8080
	router.Run(":8080")
}
````

### 自定义的验证

````
// Booking contains binded and validated data.
type Booking struct {
  CheckIn  time.Time `form:"check_in" binding:"required,bookabledate" time_format:"2006-01-02"`
  CheckOut time.Time `form:"check_out" binding:"required,gtfield=CheckIn" time_format:"2006-01-02"`
}

var bookableDate validator.Func = func(fl validator.FieldLevel) bool {
  date, ok := fl.Field().Interface().(time.Time)
  if ok {
    today := time.Now()
    if today.After(date) {
      return false
    }
  }
  return true
}

func main() {
  route := gin.Default()

  if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
    v.RegisterValidation("bookabledate", bookableDate)
  }

  route.GET("/bookable", getBookable)
  route.Run(":8085")
}

func getBookable(c *gin.Context) {
  var b Booking
  if err := c.ShouldBindWith(&b, binding.Query); err == nil {
    c.JSON(http.StatusOK, gin.H{"message": "Booking dates are valid!"})
  } else {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
  }
}

$ curl "localhost:8085/bookable?check_in=2030-04-16&check_out=2030-04-17"
{"message":"Booking dates are valid!"}

$ curl "localhost:8085/bookable?check_in=2030-03-10&check_out=2030-03-09"
{"error":"Key: 'Booking.CheckOut' Error:Field validation for 'CheckOut' failed on the 'gtfield' tag"}

$ curl "localhost:8085/bookable?check_in=2000-03-09&check_out=2000-03-10"
{"error":"Key: 'Booking.CheckIn' Error:Field validation for 'CheckIn' failed on the 'bookabledate' tag"}%
````

### bind query string

````
package main

import (
  "log"
  "net/http"

  "github.com/gin-gonic/gin"
)

type Person struct {
  Name    string `form:"name"`
  Address string `form:"address"`
}

func main() {
  route := gin.Default()
  route.Any("/testing", startPage)
  route.Run(":8085")
}

func startPage(c *gin.Context) {
  var person Person
  if c.ShouldBindQuery(&person) == nil {
    log.Println("====== Only Bind By Query String ======")
    log.Println(person.Name)
    log.Println(person.Address)
  }
  c.String(http.StatusOK, "Success")
}

````

### binding query string or post data

````
package main

import (
  "log"
  "net/http"
  "time"

  "github.com/gin-gonic/gin"
)

type Person struct {
  Name       string    `form:"name"`
  Address    string    `form:"address"`
  Birthday   time.Time `form:"birthday" time_format:"2006-01-02" time_utc:"1"`
  CreateTime time.Time `form:"createTime" time_format:"unixNano"`
  UnixTime   time.Time `form:"unixTime" time_format:"unix"`
}

func main() {
  route := gin.Default()
  route.GET("/testing", startPage)
  route.Run(":8085")
}

func startPage(c *gin.Context) {
  var person Person
  // If `GET`, only `Form` binding engine (`query`) used.
  // If `POST`, first checks the `content-type` for `JSON` or `XML`, then uses `Form` (`form-data`).
  // See more at https://github.com/gin-gonic/gin/blob/master/binding/binding.go#L88
  if c.ShouldBind(&person) == nil {
    log.Println(person.Name)
    log.Println(person.Address)
    log.Println(person.Birthday)
    log.Println(person.CreateTime)
    log.Println(person.UnixTime)
  }

  c.String(http.StatusOK, "Success")
}
````

### bing uri

```
type Person struct {
  ID string `uri:"id" binding:"required,uuid"`
  Name string `uri:"name" binding:"required"`
}

func main() {
  route := gin.Default()
  route.GET("/:name/:id", func(c *gin.Context) {
    var person Person
    if err := c.ShouldBindUri(&person); err != nil {
      c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
      return
    }
    c.JSON(http.StatusOK, gin.H{"name": person.Name, "uuid": person.ID})
  })
  route.Run(":8088")
}
```

### bind header

```
type testHeader struct {
  Rate   int    `header:"Rate"`
  Domain string `header:"Domain"`
}

func main() {
  r := gin.Default()
  r.GET("/", func(c *gin.Context) {
    h := testHeader{}

    if err := c.ShouldBindHeader(&h); err != nil {
      c.JSON(http.StatusOK, err)
    }

    fmt.Printf("%#v\n", h)
    c.JSON(http.StatusOK, gin.H{"Rate": h.Rate, "Domain": h.Domain})
  })

  r.Run()

// client
// curl -H "rate:300" -H "domain:music" 127.0.0.1:8080/
// output
// {"Domain":"music","Rate":300}
}
```

### Multipart/Urlencoded binding

```
type ProfileForm struct {
  Name   string                `form:"name" binding:"required"`
  Avatar *multipart.FileHeader `form:"avatar" binding:"required"`

  // or for multiple files
  // Avatars []*multipart.FileHeader `form:"avatar" binding:"required"`
}
func main() {
  router := gin.Default()
  router.POST("/profile", func(c *gin.Context) {
    // you can bind multipart form with explicit binding declaration:
    // c.ShouldBindWith(&form, binding.Form)
    // or you can simply use autobinding with ShouldBind method:
    var form ProfileForm
    // in this case proper binding will be automatically selected
    if err := c.ShouldBind(&form); err != nil {
      c.String(http.StatusBadRequest, "bad request")
      return
    }

    err := c.SaveUploadedFile(form.Avatar, form.Avatar.Filename)
    if err != nil {
      c.String(http.StatusInternalServerError, "unknown error")
      return
    }

    // db.Save(&form)

    c.String(http.StatusOK, "ok")
  })
  router.Run(":8080")
}
```

### XML, JSON, YAML, TOML and ProtoBuf rendering

```
func main() {
  r := gin.Default()

  // gin.H is a shortcut for map[string]any
  r.GET("/someJSON", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
  })

  r.GET("/moreJSON", func(c *gin.Context) {
    // You also can use a struct
    var msg struct {
      Name    string `json:"user"`
      Message string
      Number  int
    }
    msg.Name = "Lena"
    msg.Message = "hey"
    msg.Number = 123
    // Note that msg.Name becomes "user" in the JSON
    // Will output  :   {"user": "Lena", "Message": "hey", "Number": 123}
    c.JSON(http.StatusOK, msg)
  })

  r.GET("/someXML", func(c *gin.Context) {
    c.XML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
  })

  r.GET("/someYAML", func(c *gin.Context) {
    c.YAML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
  })

  r.GET("/someTOML", func(c *gin.Context) {
    c.TOML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
  })

  r.GET("/someProtoBuf", func(c *gin.Context) {
    reps := []int64{int64(1), int64(2)}
    label := "test"
    // The specific definition of protobuf is written in the testdata/protoexample file.
    data := &protoexample.Test{
      Label: &label,
      Reps:  reps,
    }
    // Note that data becomes binary data in the response
    // Will output protoexample.Test protobuf serialized data
    c.ProtoBuf(http.StatusOK, data)
  })

  // Listen and serve on 0.0.0.0:8080
  r.Run(":8080")
}
```

### asciiJson

asciiJSON只用来生成涉及非ascii的字符的json

```
func main() {
  r := gin.Default()

  r.GET("/someJSON", func(c *gin.Context) {
    data := gin.H{
      "lang": "GO语言",
      "tag":  "<br>",
    }

    // will output : {"lang":"GO\u8bed\u8a00","tag":"\u003cbr\u003e"}
    c.AsciiJSON(http.StatusOK, data)
  })

  // Listen and serve on 0.0.0.0:8080
  r.Run(":8080")
}
```

### save data from file

```
func main() {
  router := gin.Default()

  router.GET("/local/file", func(c *gin.Context) {
    c.File("local/file.go")
  })

  var fs http.FileSystem = // ...
  router.GET("/fs/file", func(c *gin.Context) {
    c.FileFromFS("fs/file.go", fs)
  })
}
```

### server data from reader

````
router.GET("/someDataFromReader", func(c *gin.Context) {
    response, err := http.Get("https://raw.githubusercontent.com/gin-gonic/logo/master/color.png") // 调用Get
    if err != nil || response.StatusCode != http.StatusOK {
      c.Status(http.StatusServiceUnavailable)
      return
    }

    reader := response.Body // 获取body
    defer reader.Close()
    contentLength := response.ContentLength
    contentType := response.Header.Get("Content-Type")

    extraHeaders := map[string]string{
      "Content-Disposition": `attachment; filename="gopher.png"`,
    }

    c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
  })
````

### 重定向

```
r.GET("/test", func(c *gin.Context) {
  c.Redirect(http.StatusMovedPermanently, "http://www.google.com/")
})

// 使用hanleContext重定向
r.GET("/test", func(c *gin.Context) {
    c.Request.URL.Path = "/test2"
    r.HandleContext(c)
})
r.GET("/test2", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"hello": "world"})
})
```

### 定制中间件

```
func Logger() gin.HandlerFunc {
  return func(c *gin.Context) {
    t := time.Now()

    // Set example variable
    c.Set("example", "12345")

    // before request

    c.Next()

    // after request
    latency := time.Since(t)
    log.Print(latency)

    // access the status we are sending
    status := c.Writer.Status()
    log.Println(status)
  }
}

func main() {
  r := gin.New()
  r.Use(Logger())

  r.GET("/test", func(c *gin.Context) {
    example := c.MustGet("example").(string) // MustGet获取在当前请求处理中存储的键值对（Key-Value Pair）中指定键名的值

    // it would print: "12345"
    log.Println(example)
  })

  // Listen and serve on 0.0.0.0:8080
  r.Run(":8080")
}
```

###  定制一个认证中间件

````go
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
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM) // 捕捉退出信号

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
````

curl测试

```
curl -X POST -d '{"name": "面条", "price": 10}' 'http://127.0.0.1:8080/admin/food?name=foo&password=123456'
curl -X DELETE -d '{"name": "面条", "price": 10}' 'http://127.0.0.1:8080/admin/food?name=foo&password=123456'
```

### http.ListenAndServer

定制http运行配置参数

```
func main() {
  router := gin.Default()
  http.ListenAndServe(":8080", router)
}

func main() {
  router := gin.Default()

  s := &http.Server{
    Addr:           ":8080",
    Handler:        router,
    ReadTimeout:    10 * time.Second,
    WriteTimeout:   10 * time.Second,
    MaxHeaderBytes: 1 << 20,
  }
  s.ListenAndServe()
}
```

### 运行多个服务

```go
package main

import (
  "log"
  "net/http"
  "time"

  "github.com/gin-gonic/gin"
  "golang.org/x/sync/errgroup"
)

var (
  g errgroup.Group
)

func router01() http.Handler {
  e := gin.New()
  e.Use(gin.Recovery())
  e.GET("/", func(c *gin.Context) {
    c.JSON(
      http.StatusOK,
      gin.H{
        "code":  http.StatusOK,
        "error": "Welcome server 01",
      },
    )
  })

  return e
}

func router02() http.Handler {
  e := gin.New()
  e.Use(gin.Recovery())
  e.GET("/", func(c *gin.Context) {
    c.JSON(
      http.StatusOK,
      gin.H{
        "code":  http.StatusOK,
        "error": "Welcome server 02",
      },
    )
  })

  return e
}

func main() {
  server01 := &http.Server{
    Addr:         ":8080",
    Handler:      router01(),
    ReadTimeout:  5 * time.Second,
    WriteTimeout: 10 * time.Second,
  }

  server02 := &http.Server{
    Addr:         ":8081",
    Handler:      router02(),
    ReadTimeout:  5 * time.Second,
    WriteTimeout: 10 * time.Second,
  }

  g.Go(func() error {
    err := server01.ListenAndServe()
    if err != nil && err != http.ErrServerClosed { // 判断err是否为http.ErrServerClosed
      log.Fatal(err)
    }
    return err
  })

  g.Go(func() error {
    err := server02.ListenAndServe()
    if err != nil && err != http.ErrServerClosed {
      log.Fatal(err)
    }
    return err
  })

  if err := g.Wait(); err != nil {
    log.Fatal(err)
  }
}
```

### 优雅退出

```
func main() {
  router := gin.Default()
  router.GET("/", func(c *gin.Context) {
    time.Sleep(5 * time.Second)
    c.String(http.StatusOK, "Welcome Gin Server")
  })

  srv := &http.Server{
    Addr:    ":8080",
    Handler: router,
  }

  // Initializing the server in a goroutine so that
  // it won't block the graceful shutdown handling below
  go func() {
    if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
      log.Printf("listen: %s\n", err)
    }
  }()

  // Wait for interrupt signal to gracefully shutdown the server with
  // a timeout of 5 seconds.
  quit := make(chan os.Signal)
  // kill (no param) default send syscall.SIGTERM
  // kill -2 is syscall.SIGINT
  // kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
  signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
  <-quit
  log.Println("Shutting down server...")

  // The context is used to inform the server it has 5 seconds to finish
  // the request it is currently handling
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()

  if err := srv.Shutdown(ctx); err != nil {
    log.Fatal("Server forced to shutdown:", err)
  }

  log.Println("Server exiting")
}
```

### Bind函数

```
type StructA struct {
    FieldA string `form:"field_a"`
}

type StructB struct {
    NestedStruct StructA
    FieldB string `form:"field_b"`
}

type StructC struct {
    NestedStructPointer *StructA
    FieldC string `form:"field_c"`
}

func GetDataB(c *gin.Context) {
    var b StructB
    c.Bind(&b)
    c.JSON(http.StatusOK, gin.H{
        "a": b.NestedStruct,
        "b": b.FieldB,
    })
}

func GetDataC(c *gin.Context) {
    var b StructC
    c.Bind(&b)
    c.JSON(http.StatusOK, gin.H{
        "a": b.NestedStructPointer,
        "c": b.FieldC,
    })
}

func GetDataD(c *gin.Context) {
    var b StructD
    c.Bind(&b)
    c.JSON(http.StatusOK, gin.H{
        "x": b.NestedAnonyStruct,
        "d": b.FieldD,
    })
}

func main() {
    r := gin.Default()
    r.GET("/getb", GetDataB)
    r.GET("/getc", GetDataC)
    r.GET("/getd", GetDataD)

    r.Run()
}
// curl命令
$ curl "http://localhost:8080/getb?field_a=hello&field_b=world"
{"a":{"FieldA":"hello"},"b":"world"}
$ curl "http://localhost:8080/getc?field_a=hello&field_c=world"
{"a":{"FieldA":"hello"},"c":"world"}
$ curl "http://localhost:8080/getd?field_x=hello&field_d=world"
{"d":"world","x":{"FieldX":"hello"}}
```

### 绑定body到不同的结构体中

1. `c.ShouldBindBodyWith` 可以将请求体绑定到上下文中

```
type formA struct {
  Foo string `json:"foo" xml:"foo" binding:"required"`
}

type formB struct {
  Bar string `json:"bar" xml:"bar" binding:"required"`
}

func SomeHandler(c *gin.Context) {
  objA := formA{}
  objB := formB{}
  // shouldBind消费body , body不能被复用
  if errA := c.ShouldBind(&objA); errA == nil {
    c.String(http.StatusOK, `the body should be formA`)
  // Always an error is occurred by this because c.Request.Body is EOF now.
  } else if errB := c.ShouldBind(&objB); errB == nil {
    c.String(http.StatusOK, `the body should be formB`)
  } else {
    ...
  }
}
```

### Define format for the log of routes

```
gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
    log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
}
```

### 获取或设置cookie

```
func main() {
  router := gin.Default()

  router.GET("/cookie", func(c *gin.Context) {

      cookie, err := c.Cookie("gin_cookie")

      if err != nil {
          cookie = "NotSet"
          c.SetCookie("gin_cookie", "test", 3600, "/", "localhost", false, true)
      }

      fmt.Printf("Cookie value: %s \n", cookie)
  })

  router.Run()
}
```

### 测试

The `net/http/httptest` package is preferable way for HTTP testing.

```
func TestPingRoute(t *testing.T) {
  router := setupRouter()

  w := httptest.NewRecorder()
  req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
  router.ServeHTTP(w, req)

  assert.Equal(t, http.StatusOK, w.Code)
  assert.Equal(t, "pong", w.Body.String())
}
```

