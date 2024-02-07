package controller

import (
	"fmt"
	"net/url"
	"time"

	"odisk/common"
	g "odisk/global"

	"github.com/gin-gonic/gin"
)

const (
	// Expiry for upload URL
	UploadExpiry = time.Second * 24 * 60 * 60 // 1 day.
	// Expiry for download URL
	DownloadExpiry = time.Second * 24 * 60 * 60 * 7 // 7 days.
)

// UploadFile generates a pre-signed URL for uploading a file.
// POST /object/upload
func UploadFile(c *gin.Context) {

	objectname := c.PostForm("filename")

	presignedURL, err := g.MinioClient.PresignedPutObject(g.Config.Minio.BucketName, objectname, UploadExpiry)
	if err != nil {
		common.Error(c, "生成预签名URL失败", err)
	}
	common.Success(c, "Successlly generated presigned URL", presignedURL)
}

// DownloadFile generates a pre-signed URL for downloading a file.
// GET /object/download
func DownloadFile(c *gin.Context) {
	objectname := c.PostForm("filename")

	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", objectname))
	presignedURL, err := g.MinioClient.PresignedGetObject(g.Config.Minio.BucketName, objectname, DownloadExpiry, reqParams)
	if err != nil {
		common.Error(c, "生成预签名URL失败", err)
	}
	common.Success(c, "Successlly generated presigned URL", presignedURL)
}

// DELATE /object/delate
func DelFile(c *gin.Context) {

}

// POST /object/move
func MoveFile(c *gin.Context) {

}

// POST /object/rename
func RenameFile(c *gin.Context) {

}

