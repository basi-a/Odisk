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
	mode := Config.Server.Mode
	setMode(mode)
	
	r := gin.Default()
	if mode == "debug" {
		// default is "debug/pprof"
		pprof.Register(r)
	}

	// set trusted proxys
	err := r.SetTrustedProxies(trusted_proxies)
	if err != nil {
		log.Println("SetTrustedProxies error:",err)
	}
	r.Use(sessions.Sessions("session", Store))
	// Resolve cross-domain
	r.Use(cors.Default())


	pingGroup := r.Group("/")
	{	
		
		// 使用路由组处理 HEAD 和 GET 请求
		pingGroup.HEAD("/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong\n")
		})
		pingGroup.GET("/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong\n")
		})
	}

	v1 := r.Group("/v1")
	{
		v1.GET("/hello", func(c *gin.Context) {
			common.Success(c, "helloworld", nil)
		})

		// auth := v1.Group("/", middleware.SessionMiddleware())
		// {

		// }

	}
	
	RouterEngine = r
}


func setMode(mode string) {
		
        //设置运行模式
        switch mode {
        case "debug":
                gin.SetMode(gin.DebugMode)
        case "release":
                gin.SetMode(gin.ReleaseMode)
        default:
                log.Fatalln("Your run mode is set to", mode, ". Must be debug or release!!!!")
        }
}