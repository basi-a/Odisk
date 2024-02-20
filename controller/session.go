package controller

import (
	"log"
	"net/http"


	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SaveSession(c *gin.Context, key string, value interface{}) {  
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
  
	// 将键值对存储到 session 中  
	session.Set(key, value)  
  
	// 保存 session  
	if err := session.Save(); err != nil {  
		log.Println("save err", err)  
	}  
}

func ReadSession(c *gin.Context, key string) (value interface{}) {  
	session := sessions.Default(c)  
	// 从 session 中获取值  
	value = session.Get(key)  
	return value
}

func DelSession(c *gin.Context, key string)  {
	session := sessions.Default(c)
	session.Delete(key)
}
