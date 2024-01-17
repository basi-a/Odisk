package global

import (
	"log"
	"os"

	"github.com/spf13/viper"
)
//server config
type ServerConfig struct {
	//running mode, debug or release
	Mode			string		`yaml:"mode"`
	//server runing port
	Port			string		`yaml:"port"`
	//secret
	Secret 			string		`yaml:"secret"`
	TrustedProxies 	[]string 	`yaml:"trusted_proxies"`
}
//redis config
type RedisConfig struct {
	RedisAddr		string	`yaml:"redisAddr"`
	RedisPassword	string	`yaml:"redisPassword"`
	RedisPort		string	`yaml:"redisPort"`
	RedisPoolConns	int		`yaml:"redisPoolConns"`
}
//mariadb config
type MariadbConfig struct {
	DBusername		string	`yaml:"dbUsername"`
	DBpassword 		string	`yaml:"dbPassword"`
	DBhost			string	`yaml:"dbHost"`
	DBport			string	`yaml:"dbPort"`
	DBname 			string	`yaml:"dbName"`
	Timeout			string	`yaml:"timeout"`
	DBPoolConns		int		`yaml:"dbPoolConns"`
}

//minio config
type MinioConfig struct {
	Endpoint		string	`yaml:"endpoint"`
	AccessKeyID		string	`yaml:"accessKeyID"`
	SecretAccessKey string	`yaml:"secretAccessKey"`
}
type AppConfig struct {
	// //running mode, debug or release
	// Mode			string	`yaml:"mode"`
	// //server runing port
	// Port			string	`yaml:"port"`
	// //secret
	// Secret 			string	`yaml:"secret"`

	// //redis config
	// RedisAddr		string	`yaml:"redisAddr"`
	// RedisPassword	string	`yaml:"redisPassword"`
	// RedisPort		string	`yaml:"redisPort"`
	// RedisPoolConns	int		`yaml:"redisPoolConns"`
	
	// //mariadb config
	// DBusername		string	`yaml:"dbUsername"`
	// DBpassword 		string	`yaml:"dbPassword"`
	// DBhost			string	`yaml:"dbHost"`
	// DBport			string	`yaml:"dbPort"`
	// DBname 			string	`yaml:"dbName"`
	// Timeout			string	`yaml:"timeout"`
	// DBPoolConns		int		`yaml:"dbPoolConns"`

	// //minio
	// Endpoint		string	`yaml:"endpoint"`
	// AccessKeyID		string	`yaml:"accessKeyID"`
	// SecretAccessKey string	`yaml:"secretAccessKey"`
	Server 		ServerConfig 	`yaml:"server"`
	Redis		RedisConfig		`yaml:"redis"`
	Mariadb		MariadbConfig	`yaml:"mariadb"`
	Minio 		MinioConfig		`yaml:"minio"`
}

var Paths []string = []string{
	 	"/etc/odisk/config.yaml",  
 		"/usr/local/etc/odisk/config.yaml",  
 		"./conf/config-example.yaml",  
}


// func (conf *Conf)GetConfig() *Conf {
// 	v := viper.New()

  
//  	// 遍历路径列表，尝试找到存在的配置文件  
//  	var configFile string  
//  	for _, path := range Paths {  
//  		if _, err := os.Stat(path); err == nil {  
//  			configFile = path  
//  		break  
//  		}  
//  	}
// 	v.SetConfigFile(configFile)
// 	v.SetConfigType("yaml")

// 	if err := v.ReadInConfig(); err != nil{
// 		log.Fatalln(err)
// 	}

// 	if err := v.Unmarshal(&conf); err != nil{
// 		log.Fatalln(err)
// 	}
// 	return conf
// }
var Config AppConfig
func InitConfig()  {
	// 遍历路径列表，尝试找到存在的配置文件  
 	var configFile string  
 	for _, path := range Paths {  
 		if _, err := os.Stat(path); err == nil {  
 			configFile = path  
 		break  
 		}  
 	}
	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")
	// if err := viper.ReadInConfig(); err != nil {
	// 	log.Fatalf("Error reading config file: %v", err)
	// }

	// // 将配置文件中的值写入 AppConfig 结构体
	// if err := viper.Unmarshal(&Config); err != nil {
	// 	log.Fatalf("Error unmarshalling config: %v", err)
	// }
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	if err := viper.Unmarshal(&Config); err != nil {
		log.Fatalf("Error unmarshalling config: %v", err)
	}
}
