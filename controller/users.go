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

	type JsonData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Code     string `json:"code"`
	}
	data := JsonData{}
	if err := c.ShouldBindJSON(&data); err != nil {
		common.Error(c, "绑定失败", err)
		return
	}

	value := ReadSession(c, "EmailVerifyCode")

	if emailData, ok := value.(g.EmailData); ok && emailData.Code == data.Code {
		DelSession(c, "EmailVerifyCode")

		user := m.Users{
			Email: data.Email,
		}
		log.Println(user)
		if err := user.GetUserByEmail(); err != nil {
			common.Error(c, "用户不存在", err)
			return
		}
		user.Password = data.Password // 为了判断是否存在时，不将原始密码赋值成数据库中的旧密码的密文，所有要判断后赋值
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

	user := m.Users{
		Email:    email,
		Password: password,
	}
	ok, err := user.VerifyAccount()
	if err != nil {
		common.Error(c, "认证失败", err)
	}
	if ok && err == nil {
		userInfo, err := m.GetUserInfo(email)
		// log.Println(userInfo)
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
	// 写入消息队列，等待邮件发送
	if err := g.ProduceMsg("email", "email", jsonData); err != nil {
		common.Error(c, "发送邮件错误", err)
	} else {
		SaveSession(c, "EmailVerifyCode", data)

		common.Success(c, "发送邮件成功, 请稍等", nil)
	}
}

// GET /v1/users/list  auth 组
func ListUsers(c *gin.Context) {
	userList, err := m.ListUser()
	if err != nil {
		common.Error(c, "获取失败", err)
	} else {
		common.Success(c, "获取成功", userList)
	}
}

// POST /v1/users/delete auth 组 
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
	}
	bucketmap := m.Bucketmap{
		UserID: user.ID,
	}
	if err := bucketmap.GetMap(); err != nil {
		common.Error(c, "获取用户与存储桶关联关系失败", err)
		return
	}
	if err := bucketmap.DeleteBucketMapWithTask(); err != nil {
		common.Error(c, "删除用户关联的桶以及任务失败", err)
		return
	}

	if err := g.DeactivateBucket(bucketmap.BucketName); err != nil {
		common.Error(c, "停用桶失败", err)
		return
	}
	common.Success(c, "删除成功", nil)
}

// POST /v1/users/update auth 组
func UpdateUser(c *gin.Context) {
	type JsonData struct {
		UserName   string `json:"username"`
		Password   string `json:"password"`
		Email      string `json:"email"`
		NewEmail   string `json:"newemail"`
		Permission string `json:"permission"`
		Code       string `json:"code"`
	}
	data := JsonData{}
	if err := c.ShouldBindJSON(&data); err != nil {
		common.Error(c, "绑定失败", err)
		return
	}
	if data.Code != "" {
		value := ReadSession(c, "EmailVerifyCode")
		if emailData, ok := value.(g.EmailData); !ok || emailData.Code != data.Code {
			DelSession(c, "EmailVerifyCode")
			common.Error(c, "验证码验证失败", nil)
			return
		} else {
			DelSession(c, "EmailVerifyCode")
		}

	}

	user := m.Users{
		Email: data.Email,
	}
	if err := user.GetUserByEmail(); err != nil {
		common.Error(c, "用户不存在", err)
		return
	}
	if data.NewEmail != "" {
		user.Email = data.NewEmail
	}
	if data.UserName != "" {
		user.UserName = data.UserName
	}
	if data.Permission != "" {
		user.Permission = data.Permission
	}
	if data.Password != "" {
		user.Password = data.Password
	} else {
		user.Password = ""
	}

	if err := user.Update(); err != nil {
		common.Error(c, "更新失败", err)
		return
	}
	if data.Code != "" {
		userInfo, _ := m.GetUserInfo(data.Email)
		SaveSession(c, "userInfo", userInfo)
	}
	common.Success(c, "更新成功", nil)
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

// GET /v1/logout
func Logout(c *gin.Context) {
	DelSession(c, "userInfo")
}
