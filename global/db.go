package global

import (
	"fmt"
	"log"
	"time"

	// "log"

	"gorm.io/driver/mysql" // mysql 数据库驱动
	"gorm.io/gorm"         // 使用gorm,操作数据库的 orm 框架
	"gorm.io/gorm/schema"
)

//全局db对象
var DB *gorm.DB

func InitGorm()  {
	username := Config.Mariadb.DBusername
	password := Config.Mariadb.DBpassword
	host := Config.Mariadb.DBhost
	port := Config.Mariadb.DBport
	name := Config.Mariadb.DBname
	timeout := Config.Mariadb.Timeout
	poolConns := Config.Mariadb.DBPoolConns
	log.Println(host+":"+port)
	// log.Println(c)
	var err error
	maxRetryCount := 5
	for retryCount := 0; retryCount < maxRetryCount; retryCount++{
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s", username, password, host, port, name, timeout)
	
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
		if err == nil {
			break
		}else {
			log.Println("mariadb error:",err)
		}
		time.Sleep(time.Second*20)
	}
	if err != nil {
		log.Fatalf("Failed to connect database after %d attempts\n", maxRetryCount)
	}
	sqlDB, _ := DB.DB()

	//设置连接池参数
	sqlDB.SetMaxOpenConns(poolConns)
	sqlDB.SetConnMaxIdleTime(20)

}

func GetDB() *gorm.DB{
	return DB
}