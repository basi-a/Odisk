package initialize

import (
	"encoding/gob"
	"log"
	g "odisk/global"
	m "odisk/model"
	r "odisk/router"
)

func Initialize() {
	log.Println("The application system is initializing")
	defer log.Println("Application initialization completed.")
	gob.Register(m.UserInfo{})
	g.InitConfig()
	g.InitTemplate()
	
	log.Println("Initializing system dependencies ...")
	
	g.InitRedis()
	g.InitMinio()
	g.InitNsq()
	g.InitGorm()
	
	m.InitModel()

	r.InitRouter()
}
