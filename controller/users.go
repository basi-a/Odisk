package controller

import (
	"fmt"
	"odisk/common"
	m "odisk/model"
	u "odisk/utils"
	g "odisk/global"
	"github.com/gin-gonic/gin"
)

// POST /v1/register 注册用户
func RegisterUser(c *gin.Context){
	username := c.PostForm("username")
	password := c.PostForm("password")
	email := c.PostForm("email")
	code := c.PostForm("code")
	maxUsernameLength := 20
	// 检查用户名长度
    if  len(username) > maxUsernameLength {
        common.Error(c, fmt.Sprintf("用户名长度必须在%d以内", maxUsernameLength), nil)
        return
    }

	// 验证邮箱
	value := ReadSession(c, "EmailVerifyCode")
	if emailData, ok := value.(m.EmailData); ok && emailData.Code == code{
		user := m.Users{}
		if username != "" && password != "" && email != "" {
			err := user.AddUser(username, password, email)
			if err != nil {
				common.Error(c, "注册失败, 请检查是否邮箱已使用, 或输入有误", err)
			}else{
				common.Success(c, fmt.Sprintf("注册成功, 用户名: %s", username))
			}
		}
	}else{
		common.Error(c, "邮箱验证失败", nil)
	}
}



// POST /v1/login 登陆
func Login(c *gin.Context)  {
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
		}else{
			SaveSession(c, "userInfo", userInfo)
			common.Success(c, fmt.Sprintf("Welcome %s", email))
		}
		
	}
}

// POST /v1/emailVerify
func EmailVerifyCode(c *gin.Context)   {
	email := c.PostForm("email")
	code := u.GenerateVerificationCode(6)
	object := "Verify your email address"
	
	data := m.EmailData{
		Email: email,
		Code: code,
	}
	if err := u.SendEmail(email, object, g.EmailTemplate, data); err != nil {
		common.Error(c, "发送邮件错误", err)
	}else{
		SaveSession(c, "EmailVerifyCode", data)
		common.Success(c, "发送邮件成功")
	}
}

// GET /v1/users  auth 组
func ListUsers(c *gin.Context)  {
	
}

// POST /v1/users/delate auth 组
func DelUser(c *gin.Context)  {
	
}

// POST /v1/users/update auth 组
func UpdateUser(c *gin.Context) {
	
}

// GET /v1/users/info auth 组
func UserInfo(c *gin.Context) {
	
}


// GET /userInfo
func GetUserInfo(c *gin.Context)  {

	value := ReadSession(c, "userInfo")
	// 尝试将读取的值断言为 m.UserInfo 类型
	if userInfo, ok := value.(m.UserInfo); ok && userInfo.Email != ""{
		common.Success(c, "获取成功", fmt.Sprintf("username: %s", userInfo.UserName))
	}else{
		common.Error(c, "获取信息失败", nil)
	}
}

// GET /logout
func Logout(c *gin.Context)  {
	DelSession(c, "userInfo")
}