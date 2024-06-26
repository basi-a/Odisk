package router

import (
	"log"
	"net/http"
	"odisk/controller"
	g "odisk/global"
	"odisk/middleware"
	"odisk/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"
)

func InitRouter() {
	log.Println("Route is initializing ...")
	defer log.Println("Route initialization completed.")
	r := engineBasicCfgOps()
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
		v1.POST("/emailVerify", controller.EmailVerifyCode)
		v1.POST("/reset", controller.ResetPassword)

		authedGroup := v1.Group("/", middleware.SessionMiddleware())
		{
			userGroup := authedGroup.Group("/users")
			{
				userGroup.GET("/list", controller.ListUsers)
				userGroup.POST("/update", controller.UpdateUser)
				userGroup.POST("/delete", controller.DelUser)

			}

			sessionGroup := authedGroup.Group("/")
			{
				sessionGroup.GET("/userInfo", controller.GetUserInfo)
				sessionGroup.GET("/logout", controller.Logout)
				s3Group := sessionGroup.Group("/s3")
				{
					uploadGroup := s3Group.Group("/upload")
					{
						uploadGroup.POST("/small", controller.UploadFile)
						uploadGroup.POST("/big/create", controller.MultipartUploadCreate)
						uploadGroup.POST("/big/finish", controller.MultipartUploadFinish)
						taskGroup := uploadGroup.Group("/task")
						{
							taskGroup.POST("/abort", controller.TaskAbort)
							taskGroup.POST("/del", controller.TaskDel)
							taskGroup.PUT("/add", controller.TaskAdd)
							taskGroup.PUT("/done", controller.TaskDone)
							taskGroup.PUT("/percent/update", controller.UpdateTaskPercent)
							taskGroup.GET("/percent/:taskID", controller.GetTaskPercent)
							taskGroup.POST("/list", controller.GetTaskList)
						}
					}
					s3Group.POST("/download", controller.DownloadFile)
					s3Group.POST("/delete", controller.DeleteFile)
					s3Group.POST("/mv", controller.MoveFile)
					s3Group.POST("/mkdir", controller.Mkdir)
					s3Group.POST("/list", controller.FileList)
					s3Group.POST("/size", controller.CurrentSize)
					mapGroup := s3Group.Group("/bucketmap")
					{
						mapGroup.POST("/del", controller.DeleteBucketMapWithTask)
						mapGroup.GET("/list", controller.GetMapList)
						mapGroup.POST("/update", controller.UpdateBucketmap)
					}
					
				}

			}

		}

	}

	g.RouterEngine = r
}

func engineBasicCfgOps() *gin.Engine{
	trusted_proxies := g.Config.Server.TrustedProxies
	mode := g.Config.Server.Mode
	gin.DisableConsoleColor()

	setMode(mode)

	r := gin.Default()
	if mode == "debug" {
		// default is "debug/pprof"
		pprof.Register(r)
		r.GET("/mime", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, utils.GetAllMime())
		})
	}

	// 设置信任代理
	err := r.SetTrustedProxies(trusted_proxies)
	if err != nil {
		log.Println("SetTrustedProxies error:", err)
	}
	r.Use(sessions.Sessions("session_id", g.Store))

	// 配置CORS中间件
	config := cors.DefaultConfig()
	config.AllowOrigins = append(config.AllowOrigins, g.Config.Server.CROS.AllowOrigins...)
	config.AllowCredentials = g.Config.Server.CROS.AllowCredentials

	// 将CORS中间件添加到路由引擎中
	r.Use(cors.New(config))
	return r
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
