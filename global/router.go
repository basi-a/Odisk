package global

import (
	"log"
	"net/http"
	"odisk/common"

	// "odisk/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"
)
var RouterEngine *gin.Engine
func InitRouter() {
	trusted_proxies := Config.Server.TrustedProxies
	r := gin.Default()

	// default is "debug/pprof"
	pprof.Register(r)
	
	// set trusted proxys
	err := r.SetTrustedProxies(trusted_proxies)
	if err != nil {
		log.Println("SetTrustedProxies error:",err)
	}
	r.Use(sessions.Sessions("session", Store))
	// Resolve cross-domain
	r.Use(cors.Default())

	// r.Any("/ping", func(c *gin.Context) {
	// 	if c.Request.Method == http.MethodHead || c.Request.Method == http.MethodGet {
	// 		// 返回 200 OK
	// 		c.String(http.StatusOK, "pong\n")
	// 	} else {
	// 		// 对于其他请求方法，返回 405 Method Not Allowed
	// 		c.Status(http.StatusMethodNotAllowed)
	// 	}
	// })

	pingGroup := r.Group("/ping")
	{
		// 使用路由组处理 HEAD 和 GET 请求
		pingGroup.HEAD("/", func(c *gin.Context) {
			c.String(http.StatusOK, "pong\n")
		})

		pingGroup.GET("/", func(c *gin.Context) {
			c.String(http.StatusOK, "pong\n")
		})
	}

	
	
	r.GET("/hello", func(c *gin.Context) {
		common.Success(c, "helloworld", nil)
	})

	// auth := r.Group("/", middleware.SessionMiddleware())
	// {
		
	// }

	RouterEngine = r
}