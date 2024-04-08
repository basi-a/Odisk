package controller

import (
	"crypto/md5"
	"fmt"
	"odisk/common"
	g "odisk/global"
	m "odisk/model"
	u "odisk/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// POST /v1/register 注册用户
func RegisterUser(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	email := c.PostForm("email")
	code := c.PostForm("code")
	maxUsernameLength := 20
	// 检查用户名长度
	if len(username) > maxUsernameLength {
		common.Error(c, fmt.Sprintf("用户名长度必须在%d以内", maxUsernameLength), nil)
		return
	}

	// 验证邮箱
	value := ReadSession(c, "EmailVerifyCode")
	// 生成存储桶名字
	// bucketName := base64.RawStdEncoding.EncodeToString([]byte(username + email + code))
	str := username + email + time.Now().String()
	bucketName := md5.New().Sum([]byte(str))
	if emailData, ok := value.(m.EmailData); ok && emailData.Code == code {
		user := m.Users{}
		if username != "" && password != "" && email != "" {

			userID, err := user.AddUser(username, password, email, nil)
			bucketmap := m.Bucketmap{
				UserID:     *userID,
				BucketName: string(bucketName),
			}
			if err != nil {
				common.Error(c, "注册失败, 请检查是否邮箱已使用, 或输入有误", err)
			} else if err := bucketmap.SaveMap(); err != nil {
				common.Error(c, "用户与存储桶对应关系记录失败", err)
			} else {
				common.Success(c, fmt.Sprintf("注册成功, 用户名: %s", username))
			}

			// 创建存储桶
			// log.Println(bucketName)
			if err := g.MakeBucket(bucketmap.BucketName); err != nil {
				common.Error(c, "创建存储桶失败", err)
			} else {
				common.Success(c, "存储桶创建成功")
			}
		}

	} else {
		common.Error(c, "邮箱验证失败", nil)
	}

}

// POST /v1/login 登陆
func Login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	user := m.Users{}
	ok, err := user.VerifyAccount(email, password)
	if err != nil {
		common.Error(c, "认证失败", err)
	}
	if ok && err == nil {
		userInfo, err := m.GetUserInfo(email)
		if err != nil {
			common.Error(c, "获取信息失败", err)
		} else {
			SaveSession(c, "userInfo", userInfo)
			common.Success(c, fmt.Sprintf("Welcome %s", email))
		}

	}
}

// POST /v1/emailVerify
func EmailVerifyCode(c *gin.Context) {
	email := c.PostForm("email")
	code := u.GenerateVerificationCode(6)
	subject := "Verify your email address"

	data := m.EmailData{
		Email: email,
		Code:  code,
	}
	nsqBody := g.EmailData{
		Email: 		email,
		Subject: 	subject,
		T:			g.EmailTemplate,
		Data: 		data,
	}
	if err := g.NsqPublish("email", nsqBody); err != nil {
		common.Error(c, "发送邮件错误", err)
	} else {
		SaveSession(c, "EmailVerifyCode", data)
		common.Success(c, "发送邮件成功, 请稍等")
	}
}

// GET /v1/users  auth 组
func ListUsers(c *gin.Context) {

}

// POST /v1/users/delate auth 组
func DelUser(c *gin.Context) {

}

// POST /v1/users/update auth 组
func UpdateUser(c *gin.Context) {

}

// GET /v1/users/info auth 组
func UserInfo(c *gin.Context) {

}

// GET /userInfo
func GetUserInfo(c *gin.Context) {

	value := ReadSession(c, "userInfo")
	// 尝试将读取的值断言为 m.UserInfo 类型
	if userInfo, ok := value.(m.UserInfo); ok && userInfo.Email != "" {
		common.Success(c, "获取成功", fmt.Sprintf("username: %s", userInfo.UserName))
	} else {
		common.Error(c, "获取信息失败", nil)
	}
}

// GET /logout
func Logout(c *gin.Context) {
	DelSession(c, "userInfo")
}
