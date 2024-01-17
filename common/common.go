package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, msg string, data interface{})  {
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":	msg,
		"data":	data,
	})
}

func Error(c *gin.Context, msg string)  {
	c.JSON(http.StatusOK, gin.H{
		"error": msg,
	})
}