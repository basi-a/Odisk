package global

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

// server config
type ServerConfig struct {
	//running mode, debug or release
	Mode string `yaml:"mode"`
	//server runing port
	Port string `yaml:"port"`
	//secret
	Secret         string      `yaml:"secret"`
	Ssl            SslConfig   `yaml:"ssl"`
	TrustedProxies []string    `yaml:"trusted_proxies"`
	Mail           MailConfig  `yaml:"mail"`
	CROS           CrosConfig  `yaml:"cros"`
	Admin          AdminConfig `yaml:"admin"`
}

type AdminConfig struct {
	UserAdmin UserAndS3AdminConfig `yaml:"useradmin"`
	S3Admin   UserAndS3AdminConfig `yaml:"s3admin"`
}
type UserAndS3AdminConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Email    string `yaml:"email"`
}

type CrosConfig struct {
	AllowOrigins     []string `yaml:"allowOrigins"`
	AllowCredentials bool     `yaml:"allowCredentials"`
}

type MailConfig struct {
	SenderMail string `yaml:"senderMail"`
	SmtpServer string `yaml:"smtpServer"`
	Port       int    `yaml:"port"`
	UserName   string `yaml:"username"`
	Password   string `yaml:"password"`
	Template   string `yaml:"template"`
}

type SslConfig struct {
	PrivateKey string `yaml:"privatekey"`
	Cert       string `yaml:"cert"`
}

// redis config
type RedisConfig struct {
	RedisAddr      string `yaml:"redisAddr"`
	RedisPassword  string `yaml:"redisPassword"`
	RedisPort      string `yaml:"redisPort"`
	RedisPoolConns int    `yaml:"redisPoolConns"`
}
type DatabaseConfig struct {
	Dbselect string        `yaml:"dbselect"`
	Mariadb  MariadbConfig `yaml:"mariadb"`
	Pgsql    PgsqlConfig   `yaml:"pgsql"`
}

// mariadb config
type MariadbConfig struct {
	DBusername  string `yaml:"dbUsername"`
	DBpassword  string `yaml:"dbPassword"`
	DBhost      string `yaml:"dbHost"`
	DBport      string `yaml:"dbPort"`
	DBname      string `yaml:"dbName"`
	Timeout     string `yaml:"timeout"`
	DBPoolConns int    `yaml:"dbPoolConns"`
}

type PgsqlConfig struct {
	//host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai
	DBusername  string `yaml:"dbUsername"`
	DBpassword  string `yaml:"dbPassword"`
	DBhost      string `yaml:"dbHost"`
	DBport      string `yaml:"dbPort"`
	DBname      string `yaml:"dbName"`
	Sslmode     string `yaml:"sslmode"`
	TimeZone    string `yaml:"timezone"`
	DBPoolConns int    `yaml:"dbPoolConns"`
}

// minio config
type MinioConfig struct {
	Endpoint        string `yaml:"endpoint"`
	AccessKeyID     string `yaml:"accessKeyID"`
	SecretAccessKey string `yaml:"secretAccessKey"`
	UseSSL          bool   `yaml:"usessl"`
	Location        string `yaml:"location"`
	BucketMaxSize   int    `yaml:"bucketMaxSize"`
}
type NsqConfig struct {
	Port struct {
		Nsqlookupd struct {
			TCP  string `mapstructure:"tcp"`
			HTTP string `mapstructure:"http"`
		} `mapstructure:"nsqlookupd"`
		Nsqd struct {
			TCP  string `mapstructure:"tcp"`
			HTTP string `mapstructure:"http"`
		} `mapstructure:"nsqd"`
	} `mapstructure:"port"`
	Nsqlookupd []string            `yaml:"nsqlookupd"`
	Nsqd       string              `yaml:"nsqd"`
	Topics     map[string][]string `mapstructure:"topics"`
}

type AppConfig struct {
	Server   ServerConfig   `yaml:"server"`
	Redis    RedisConfig    `yaml:"redis"`
	Database DatabaseConfig `yaml:"database"`
	Minio    MinioConfig    `yaml:"minio"`
	Nsq      NsqConfig      `yaml:"nsq" mapstructure:"nsq"`
}

var Paths []string = []string{
	"/etc/odisk/config.yaml",
	"/usr/local/etc/odisk/config.yaml",
	"./conf/config-example.yaml",
}

var Config AppConfig

func InitConfig() {
	log.Println("Reading global configuration file ...")
	defer log.Println("Global configuration reading completed.")
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

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	if err := viper.Unmarshal(&Config); err != nil {
		log.Fatalf("Error unmarshalling config: %v", err)
	}

}
