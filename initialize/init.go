package initialize

import(
	g "odisk/global"
) 
func Initialize()  {
	g.InitConfig()
	g.InitGorm()
	g.InitRedis()
	g.InitRouter()
}