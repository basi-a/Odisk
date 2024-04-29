package model

type UserInfo struct {
	Email            string `json:"email"`
	UserName         string `json:"username"`
	Registrationtime string `json:"registrationtime"`
	BucketName       string `json:"bucketname"`

	Permission       string `json:"permission"`
}

type FileInfo struct {
	Key          string `json:"name"`         // Name of the object
	LastModified string `json:"lastModified"` // Date and time the object was last modified.
	Size         string `json:"size"`         // Size in bytes of the object.
	IsDir        bool   `json:"isdir"`
	ContentType  string `json:"contenttype"`
}

func GetUserInfo(email string) (userInfo UserInfo, err error) {

	user := Users{
		Email: email,
	}

	if err := user.GetUserByEmail(); err != nil {
		return UserInfo{}, err
	}

	bucketmap := Bucketmap{
		UserID: user.ID,
	}
	if err := bucketmap.GetUserBucketName(); err != nil {
		return UserInfo{}, err
	}
	userInfo = UserInfo{
		Email:            user.Email,
		UserName:         user.UserName,
		Registrationtime: user.CreatedAt.Format("2006-01-02 15:04:05"),
		BucketName:       bucketmap.BucketName,
		// BucketmapID:      bucketmap.ID,
		Permission:       user.Permission,
	}

	return
}
