package common

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  msg,
		"data": data,
	})
}

func Error(c *gin.Context, msg string, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"code":  http.StatusBadRequest,
		"msg":   msg,
		"error": err.Error(),
	})
	log.Printf("An error has occurred: %s\n", err.Error())
}
