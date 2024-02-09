package controller

import (
	"fmt"
	"log"
	"odisk/common"
	m "odisk/model"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SaveSession(c *gin.Context, userInfo m.UserInfo)  {
	session := sessions.Default(c)
	session.Options(sessions.Options{
		MaxAge: int(time.Hour*24),
		Secure: true,
		HttpOnly: true,
	})
	session.Set("userInfo", userInfo)
	if err := session.Save(); err != nil{
		log.Println("save err",err)
	}
}

// GET /userInfo
func ReadSession(c *gin.Context)  {
	session := sessions.Default(c)
	userInfo := session.Get("userInfo").(m.UserInfo)
	if userInfo.Email != ""{
		common.Success(c, "获取成功", fmt.Sprintf("username: %s", userInfo.UserName))
	}else{
		common.Error(c, "获取信息失败", nil)
	}
}

// GET /logout
func Logout(c *gin.Context)  {
	session := sessions.Default(c)
	session.Delete("userInfo")
}