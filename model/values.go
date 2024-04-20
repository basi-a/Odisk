package model

type UserInfo struct {
	Email      string `json:"email"`
	UserName   string `json:"username"`
	Permission string `json:"permission"`
}

func GetUserInfo(email string) (userInfo UserInfo, err error) {

	var user Users

	if user, err = user.GetUser(email); err != nil {
		return UserInfo{}, err
	}

	userInfo = UserInfo{
		Email:      user.Email,
		UserName:   user.UserName,
		Permission: user.Permission,
	}

	return
}
