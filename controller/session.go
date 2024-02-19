package controller

import (
	"fmt"
	"log"
	"net/http"
	"odisk/common"
	m "odisk/model"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SaveSession(c *gin.Context, userInfo m.UserInfo)  {
	session := sessions.Default(c)

	// 设置 session 选项  
	session.Options(sessions.Options{  
		// Domain:         ".example.com",                   // 可选：设置 cookie 的域名  
		// Path:           "/",                              // 可选：设置 cookie 的路径  
		MaxAge:         86400 * 7,                        // 1 周  
		Secure:         true,                             // 强制通过 HTTPS 传输  
		HttpOnly:       true,                             // 限制 JavaScript 访问  
		SameSite:       http.SameSiteLaxMode,              // 设置 SameSite 为 Lax  
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