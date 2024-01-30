package router

import (
	"log"
	"net/http"
	"odisk/controller"
	g "odisk/global"
	"odisk/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"
)

func InitRouter() {
	trusted_proxies := g.Config.Server.TrustedProxies
	mode := g.Config.Server.Mode
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
	r.Use(sessions.Sessions("session_id", g.Store))
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
		v1.POST("/register", controller.RegisterUser)
		v1.POST("/login", controller.Login)

		authedGroup := v1.Group("/", middleware.SessionMiddleware())
		{
			userGroup := authedGroup.Group("/users")
			{
				userGroup.GET("/", controller.ListUsers)
				userGroup.GET("/info", controller.UserInfo)
				userGroup.POST("/update", controller.UpdateUser)
				userGroup.POST("/delate", controller.DelUser)
			}

			sessionGroup := authedGroup.Group("/")
			{	
				sessionGroup.GET("/userInfo", controller.ReadSession)
				sessionGroup.GET("/logout", controller.Logout)
			}

		}

	}
	
	g.RouterEngine = r
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