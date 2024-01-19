package initialize

import(
	g "odisk/global"
	m "odisk/model"
) 
func Initialize()  {
	g.InitConfig()
	g.InitGorm()
	m.InitModel()
	g.InitRedis()
	g.InitMinio()
	g.InitRouter()
}