package main

import (
	"log"
	// "net/http"
	"odisk/initialize"
	g "odisk/global"

	"github.com/gin-gonic/gin"
)

func main()  {
	
	port := g.Config.Server.Port
	setMode()
	r := g.RouterEngine
	r.Run(port) 
}

func init()  {
	initialize.Initialize()
}

func setMode() {
		mode := g.Config.Server.Mode
		
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