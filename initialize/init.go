package initialize

import(
	g "odisk/global"
	m "odisk/model"
	r "odisk/router"
) 
func Initialize()  {
	g.InitConfig()
	g.InitGorm()
	m.InitModel()
	g.InitRedis()
	g.InitMinio()
	r.InitRouter()
}