package conf

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql" // mysql 数据库驱动
	"gorm.io/gorm"         // 使用gorm,操作数据库的 orm 框架
	"gorm.io/gorm/schema"
)

//全局db对象
var DB *gorm.DB

func InitGorm()  {
	conf := new(Conf)
	c := conf.GerConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s", c.DBusername, c.DBpassword, c.DBhost, c.DBport, c.DBname, c.Timeout)
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatalln("database connection err", err)
	}
	sqlDB, _ := DB.DB()

	//设置连接池参数
	sqlDB.SetMaxOpenConns(2048)
	sqlDB.SetConnMaxIdleTime(20)

}

func GetDB() *gorm.DB{
	return DB
}