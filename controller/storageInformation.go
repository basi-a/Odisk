package controller

import "github.com/gin-gonic/gin"

func FileList(c *gin.Context) {
	c.PostForm("number_per_page")
	c.PostForm("pages")
}
