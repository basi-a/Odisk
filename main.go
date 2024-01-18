package main

import (
	// "log"
	// "net/http"
	"odisk/initialize"
	g "odisk/global"

	// "github.com/gin-gonic/gin"
)

func main()  {
	
	port := g.Config.Server.Port
	
	r := g.RouterEngine
	
	r.Run(port) 
}

func init()  {
	initialize.Initialize()
}

