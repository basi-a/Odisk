package model

import (
	g "odisk/global"
)
type Info struct {
	Email		string  `json:"email"`
	UserName 	string `json:"username"`
}


func (info *Info)GetInfo(email string) error{
	db := g.DB
	if err := db.First(&info).Where("email = ?", email).Error; err != nil {
		return err
	}
	return nil
}