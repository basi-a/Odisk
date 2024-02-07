package model

import (
	"time"

	// "gorm.io/gorm"
)

type StorageInformation struct {
	Id			int
	FileName	string
	Uri			string
	UUID		string
	UploadTime	time.Time
}

