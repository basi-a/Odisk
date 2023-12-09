package main

import (
	"log"
	"odisk/conf"
)

func main()  {
	conf := new(conf.Conf)
	c := conf.GerConfig()
	log.Println("running mode", c.Mode)
}

func init()  {
	conf.InitGorm()
}