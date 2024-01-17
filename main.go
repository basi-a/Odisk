package main

import (
	"log"
	// "net/http"
	"odisk/conf"
	"odisk/initialize"

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
	
	
	r.Run(c.Port) 
}

func init()  {
	initialize.Initialize()
}