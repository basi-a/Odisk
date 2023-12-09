package conf

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Conf struct {
	//running mode, debug or release
	Mode			string	`yaml:"mode"`
	//server runing port
	Port			string	`yaml:"port"`
	//secret
	Secret 			string	`yaml:"secret"`

	//redis config
	RedisAddr		string	`yaml:"redisAddr"`
	RedisPassword	string	`yaml:"redisPassword"`
	RedisPort		string	`yaml:"redisPort"`
	
	//mariadb config
	DBusername		string	`yaml:"dbUsername"`
	DBpassword 		string	`yaml:"dbPassword"`
	DBhost			string	`yaml:"dbHost"`
	DBport			string	`yaml:"dbPort"`
	DBname 			string	`yaml:"dbName"`
	Timeout			string	`yaml:"timeout"`
}

const (
	path0="/etc/odisk/config.yaml"
	path1="/usr/local/etc/odisk/config.yaml"
	example="./conf/config-example.yaml"
)

func (conf *Conf)GerConfig() *Conf {
	v := viper.New()

	if _, err := os.Stat(path0); err == nil{
		v.SetConfigFile(path0)
		// log.Println("conf from",path0)
	}else if _, err := os.Stat(path1); err == nil{
		v.SetConfigFile(path1)
		// log.Println("conf from", path1)
	}else{
		v.SetConfigFile(example)
		// log.Println("conf from",example)
	}

	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil{
		log.Fatalln(err)
	}

	if err := v.Unmarshal(&conf); err != nil{
		log.Fatalln(err)
	}
	return conf
}