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

	//minio
	Endpoint		string	`yaml:"endpoint"`
	AccessKeyID		string	`yaml:"accessKeyID"`
	SecretAccessKey string	`yaml:"secretAccessKey"`
}

// const (
// 	path0="/etc/odisk/config.yaml"
// 	path1="/usr/local/etc/odisk/config.yaml"
// 	example="./conf/config-example.yaml"
// )

var Paths []string = []string{
	 	"/etc/odisk/config.yaml",  
 		"/usr/local/etc/odisk/config.yaml",  
 		"./conf/config-example.yaml",  
}


func (conf *Conf)GetConfig() *Conf {
	v := viper.New()

  
 	// 遍历路径列表，尝试找到存在的配置文件  
 	var configFile string  
 	for _, path := range Paths {  
 		if _, err := os.Stat(path); err == nil {  
 			configFile = path  
 		break  
 		}  
 	}
	v.SetConfigFile(configFile)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil{
		log.Fatalln(err)
	}

	if err := v.Unmarshal(&conf); err != nil{
		log.Fatalln(err)
	}
	return conf
}