package model

import (
	g "odisk/global"

	"gorm.io/gorm"
)

type Bucketmap struct {
	gorm.Model
	UserID     uint `gorm:"uniqueIndex"`
	BucketName string
}


func AutoMigrateBucketmap() {
	if g.DB.Migrator().HasTable(&Bucketmap{}) {
		return
	}
	g.DB.AutoMigrate(&Bucketmap{})

}
func (bucketmap *Bucketmap) SaveMap() error {
	db := g.DB
	if err := db.Create(&bucketmap).Error; err != nil {
		return err
	}
	return nil
}

func (bucketmap *Bucketmap) GetUserBucketName(id uint)error  {
	db := g.DB
	if err := db.Select("bucket_name").Where("user_id=?", id).First(&bucketmap).Error; err != nil {
		return err
	}
	return nil
}