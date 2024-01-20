package controller

import (
	"fmt"
	"odisk/common"
	m "odisk/model"
	u "odisk/utils"

	// "github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// POST /v1/register 注册用户
func RegisterUser(c *gin.Context){
	username := c.PostForm("username")
	password := c.PostForm("password")
	email := c.PostForm("email")
	maxUsernameLength := 20
	// 检查用户名长度
    if  len(username) > maxUsernameLength {
        common.Error(c, fmt.Sprintf("用户名长度必须在%d以内", maxUsernameLength))
        return
    }

    // 检查邮箱格式
    if !u.IsValidEmail(email) {
        common.Error(c, "邮箱格式不正确")
        return
    }
	user := m.Users{}
	if username != "" && password != "" && email != "" {
		err := user.AddUser(username, password, email)
		if err != nil {
			common.Error(c, "注册失败, 请检查是否邮箱已使用, 或输入有误")
		}else{
			common.Success(c, fmt.Sprintf("注册成功, 用户名: %s", username),nil)
		}
	}
}



// POST /v1/login 登陆
func Login(c *gin.Context)  {
	email := c.PostForm("email")
	password := c.PostForm("password")

	user := m.Users{}
	if user.VerifyAccount(email, password) {
		
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