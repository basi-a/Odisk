package main

import (
	"fmt"
	"odisk/conf"
)

func main()  {
	conf := new(conf.Conf)
	c := conf.GerConfig()
	fmt.Println(c.Mode)
}