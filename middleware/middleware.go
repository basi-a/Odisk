package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if session.Get("uuid") == nil {
			c.Abort()
			c.JSON(http.StatusOK, "No permissions, please log in.")
		} else {
			c.Next()
		}
	}
}