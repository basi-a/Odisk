package controller

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
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
	str := username + email + time.Now().String()
	bucketName := fmt.Sprintf("%x", md5.Sum([]byte(str)))
	log.Println(bucketName)
	if emailData, ok := value.(g.EmailData); ok && emailData.Code == code {
		DelSession(c, "EmailVerifyCode")

		if username != "" && password != "" && email != "" {
			user := m.Users{
				UserName:   username,
				Password:   password,
				Email:      email,
				Permission: "",
			}
			userID, err := user.AddUser()
			bucketmap := m.Bucketmap{
				UserID:     userID,
				BucketName: bucketName,
			}
			if err != nil {
				common.Error(c, "注册失败, 请检查是否邮箱已使用, 或输入有误", err)
				return
			} else if err := bucketmap.SaveMap(); err != nil {
				common.Error(c, "用户与存储桶对应关系记录失败", err)
				return
			} else {
				common.Success(c, fmt.Sprintf("注册成功, 用户名: %s", username), nil)
				// 创建存储桶
				if err := g.MakeBucket(bucketmap.BucketName); err != nil {
					log.Println("创建存储桶失败", err)
					return
				} else {
					log.Println("存储桶创建成功")
				}
			}

		} else {
			common.Error(c, "用户名、邮箱、密码都不能为空", nil)
			return
		}

	} else {
		common.Error(c, "邮箱验证失败", nil)
	}

}

// POST /v1/reset
func ResetPassword(c *gin.Context) {
	password := c.PostForm("password")
	email := c.PostForm("email")
	code := c.PostForm("code")

	value := ReadSession(c, "EmailVerifyCode")

	if emailData, ok := value.(g.EmailData); ok && emailData.Code == code {
		DelSession(c, "EmailVerifyCode")
		user := m.Users{
			Password: password,
			Email:    email,
		}
		if err := user.GetUserByEmail(); err != nil {
			common.Error(c, "用户不存在", err)
			return
		}
		if err := user.Update(); err != nil {
			common.Error(c, "更新失败", err)
			return
		}
	}
	common.Success(c, "重置成功", nil)
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
		log.Println(userInfo)
		if err != nil {
			common.Error(c, "获取信息失败", err)
		} else {
			SaveSession(c, "userInfo", userInfo)
			common.Success(c, fmt.Sprintf("Welcome %s", email), nil)
		}

	}
}

// POST /v1/emailVerify 发送邮件验证码
func EmailVerifyCode(c *gin.Context) {
	email := c.PostForm("email")
	code := u.GenerateVerificationCode(6)
	subject := "Verify your email address"

	data := g.EmailData{
		Email: email,
		Code:  code,
	}
	jsonData, _ := json.Marshal(data)
	base64str := base64.RawStdEncoding.EncodeToString(jsonData)
	sendEmailData := g.SendEmailData{
		Email:      email,
		Subject:    subject,
		DataBase64: base64str,
	}

	jsonData, _ = json.Marshal(sendEmailData)
	if err := g.ProduceMsg("email", "email", jsonData); err != nil {
		common.Error(c, "发送邮件错误", err)
	} else {
		SaveSession(c, "EmailVerifyCode", data)
		common.Success(c, "发送邮件成功, 请稍等", nil)
	}
}

// GET /v1/users  auth 组
func ListUsers(c *gin.Context) {
	userList, err := m.ListUser()
	if err != nil {
		common.Error(c, "获取失败", err)
	} else {
		common.Success(c, "获取成功", userList)
	}
}

// POST /v1/users/delate auth 组 这个请求完必须 请求 // DELATE /s3/bucketmapdel
func DelUser(c *gin.Context) {
	type JsonData struct {
		Email string `json:"email"`
	}
	data := JsonData{}
	if err := c.ShouldBindJSON(&data); err != nil {
		common.Error(c, "绑定失败", err)
		return
	}
	user := m.Users{
		Email: data.Email,
	}
	if err := user.GetUserByEmail(); err != nil {
		common.Error(c, "用户不存在", err)
		return
	}
	if err := user.DelUser(); err != nil {
		common.Error(c, "删除失败", err)
		return
	} else {
		common.Success(c, "删除成功", err)
	}
}

// POST /v1/users/update auth 组
func UpdateUser(c *gin.Context) {
	type JsonData struct {
		UserName string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	data := JsonData{}
	if err := c.ShouldBindJSON(&data); err != nil {
		common.Error(c, "绑定失败", err)
		return
	}

	user := m.Users{
		UserName: data.UserName,
		Password: data.Password,
		Email:    data.Email,
	}
	if err := user.GetUserByEmail(); err != nil {
		common.Error(c, "用户不存在", err)
		return
	}
	if err := user.Update(); err != nil {
		common.Error(c, "更新失败", err)
	}

}

// GET /v1/userInfo
func GetUserInfo(c *gin.Context) {

	value := ReadSession(c, "userInfo")

	// 尝试将读取的值断言为 m.UserInfo 类型
	if userInfo, ok := value.(m.UserInfo); ok && userInfo.Email != "" {
		common.Success(c, "获取成功", userInfo)
	} else {
		common.Error(c, "获取信息失败", nil)
	}
}

// GET /logout
func Logout(c *gin.Context) {
	DelSession(c, "userInfo")
}
