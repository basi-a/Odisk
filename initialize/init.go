package initialize

import (
	"encoding/gob"
	g "odisk/global"
	m "odisk/model"
	r "odisk/router"
) 
func Initialize()  {
	gob.Register(m.UserInfo{})
	g.InitConfig()
	g.InitGorm()
	m.InitModel()
	g.InitRedis()
	g.InitMinio()
	r.InitRouter()
}