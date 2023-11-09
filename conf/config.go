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
	
	//mariadb config
	DBusername		string	`yaml:"dbUsername"`
	DBpassword 		string	`yaml:"dbPassword"`
	DBhost			string	`yaml:"dbHost"`
	DBport			string	`yaml:"dbPort"`
	DBname 			string	`yaml:"dbName"`
	Timeout			string	`yaml:"timeout"`
}

func (conf *Conf)GerConfig() *Conf {
	v := viper.New()
	
	if _, err := os.Stat("/etc/odisk/config.yaml"); err == nil{
		v.SetConfigFile("/etc/odisk/config.yaml")
		log.Println("conf path0")
	}else if _, err := os.Stat("/usr/local/etc/odisk/config.yaml"); err == nil{
		v.SetConfigFile("/usr/local/etc/odisk/config.yaml")
		log.Println("conf path1")
	}else{
		v.SetConfigFile("./conf/config-example.yaml")
		log.Println("conf example")
	}
	// if _, err := os.Stat("/etc/odisk/config.yaml"); os.IsNotExist(err){
	// 	v.SetConfigFile("./conf/config-example.yaml")
		
	// }else { 
	// 	v.SetConfigFile("/etc/odisk/config.yaml")
	// }
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil{
		log.Fatalln(err)
	}

	if err := v.Unmarshal(&conf); err != nil{
		log.Fatalln(err)
	}
	return conf
}