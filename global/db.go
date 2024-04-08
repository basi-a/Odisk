package global

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// 全局db对象
var DB *gorm.DB

func InitGorm() {
	
	switch Config.Database.Dbselect {
	case "mariadb":
		RetryWithExponentialBackoff(UseMysql, "Database Connection")
	case "pgsql":
		RetryWithExponentialBackoff(UsePgsql, "Database Connection")
	}
}

func GetDB() *gorm.DB {
	return DB
}

func UseMysql() error {

	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s",
		Config.Database.Mariadb.DBusername,
		Config.Database.Mariadb.DBpassword,
		Config.Database.Mariadb.DBhost,
		Config.Database.Mariadb.DBport,
		Config.Database.Mariadb.DBname,
		Config.Database.Mariadb.Timeout)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		return err
	}
	sqlDB, _ := DB.DB()

	//设置连接池参数
	sqlDB.SetMaxOpenConns(Config.Database.Mariadb.DBPoolConns)
	sqlDB.SetConnMaxIdleTime(20)

	return nil
}

func UsePgsql() error {

	var err error

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		Config.Database.Pgsql.DBhost,
		Config.Database.Pgsql.DBusername,
		Config.Database.Pgsql.DBpassword,
		Config.Database.Pgsql.DBname,
		Config.Database.Pgsql.DBport,
		Config.Database.Pgsql.Sslmode,
		Config.Database.Pgsql.TimeZone)
	DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		return err
	}
	sqlDB, _ := DB.DB()

	//设置连接池参数
	sqlDB.SetMaxOpenConns(Config.Database.Pgsql.DBPoolConns)
	sqlDB.SetConnMaxIdleTime(20)
	return nil
}
