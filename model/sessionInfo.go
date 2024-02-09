package model

import (
	g "odisk/global"
)
type UserInfo struct {
	Email		string  `json:"email"`
	UserName 	string `json:"username"`
}


func GetUserInfo(email string) (userInfo UserInfo, err error){
	db := g.DB
	// if err := db.First(&info).Where("email = ?", email).Error; err != nil {
	// 	return err
	// }
	user := Users{}
	if err := db.First(&user).Select("userName").Where("email = ?", email).Error; err != nil {
		return UserInfo{}, err
	}
	userInfo = UserInfo{
		Email: user.Email,
		UserName: user.UserName,
	}
	return
}