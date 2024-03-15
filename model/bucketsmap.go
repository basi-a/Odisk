package model

import (
	g "odisk/global"
)

type Bucketmap struct {
	UserID	uint
	BucketName string
}

func (bucketmap *Bucketmap)SaveMap() error {
	db := g.DB
	if err := db.Create(&bucketmap).Error; err != nil {
		return err
	}
	return nil
}