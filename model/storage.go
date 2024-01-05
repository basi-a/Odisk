package model

import (
	"time"

	"gorm.io/gorm"
)

type Storage struct {
	Id			int
	FileName	string
	Uri			string
	UUID		string
	UploadTime	time.Time
}

