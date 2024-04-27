package model

import "log"

func InitModel()  {
	defer log.Println("Initialization of database table completed. ")
	AutoMigrateUser()
	AutoMigrateBucketmapAndTaskList()
}