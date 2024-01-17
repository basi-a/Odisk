package main

import (
	"log"
	"odisk/conf"
	"github.com/gin-gonic/gin"
)

func main()  {
	conf := new(conf.Conf)
	c := conf.GetConfig()
	log.Println("running mode", c.Mode)
	runningMode := gin.ReleaseMode
	if c.Mode != "release" {
		runningMode = gin.DebugMode
	}
	gin.SetMode(runningMode)
	r := gin.Default()
	
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "helloworld",
		})
	})
	r.Run(c.Port) 
}

func init()  {
	conf.InitGorm()
}