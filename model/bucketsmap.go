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

func (bucketmap *Bucketmap) SaveMap() error {
	db := g.DB
	if err := db.Create(&bucketmap).Error; err != nil {
		return err
	}
	return nil
}

func InitBucketmap() {
	g.DB.AutoMigrate(Bucketmap{})
}
