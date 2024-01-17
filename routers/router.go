package routers

import (
	"fmt"
	"log"
	"net/http"
	"odisk/common"
	"odisk/conf"

	// "odisk/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	conf := new(conf.Conf)
	c := conf.GetConfig()
	
	store, err := redis.NewStore(c.RedisPoolConns, "tcp", fmt.Sprintf("%s:%s",c.RedisAddr,c.RedisPort), c.RedisPassword, []byte(c.Secret))
	if err != nil {
		log.Fatalln(err)
	}
	r.Use(sessions.Sessions("session", store))
	r.Use(cors.Default())

	r.Any("/ping", func(c *gin.Context) {
		if c.Request.Method == http.MethodHead || c.Request.Method == http.MethodGet {
			// 返回 200 OK
			c.String(http.StatusOK, "pong\n")
		} else {
			// 对于其他请求方法，返回 405 Method Not Allowed
			c.Status(http.StatusMethodNotAllowed)
		}
	})
	r.GET("/hello", func(c *gin.Context) {
		common.Success(c, "helloworld", nil)
	})

	// auth := r.Group("/", middleware.SessionMiddleware())
	// {
		
	// }

	return r
}