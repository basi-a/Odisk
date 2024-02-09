package controller

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"time"

	"odisk/common"
	g "odisk/global"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

const (
	// Expiry for upload URL
	UploadExpiry = time.Second * 24 * 60 * 60 // 1 day.
	// Expiry for download URL
	DownloadExpiry = time.Second * 24 * 60 * 60 * 7 // 7 days.
)

// UploadFile generates a pre-signed URL for uploading a file.
// POST /s3/upload/small
func UploadFile(c *gin.Context) {
	objectname := c.PostForm("objectname")

	presignedURL, err := g.S3Client.PresignedPutObject(g.S3Ctx, g.Config.Minio.BucketName, objectname, UploadExpiry)
	if err != nil {
		common.Error(c, "生成预签名URL失败", err)
		return
	} else {
		uri := fmt.Sprintf("uploadUri: %s", presignedURL.RequestURI())
		common.Success(c, "Successlly generated presigned URL", uri)
	}

}

// POST /s3/upload/big/create
func MultipartUploadCreate(c *gin.Context) {
	objectname := c.PostForm("objectname")
	partNumberArr := c.PostFormArray("partNumberArr")

	uploadID, err := g.S3Client.NewMultipartUpload(g.S3Ctx, g.Config.Minio.BucketName, objectname, minio.PutObjectOptions{})
	if err != nil {
		common.Error(c, "生成uploadID失败", err)
	}

	presignedURLs := make([]string, 0)
	for _, v := range partNumberArr {

		presignedURL, err := g.S3Client.Presign(g.S3Ctx, "PUT", g.Config.Minio.BucketName, objectname, UploadExpiry, url.Values{
			"uploadID":   []string{uploadID},
			"partNumber": []string{v},
		})
		if err != nil {
			common.Error(c, "生成预签名URL失败", err)
			return
		}
		presignedURLs = append(presignedURLs, presignedURL.RequestURI())
	}
	data := make([]string,3)
	data[0] = fmt.Sprintf("uploadID: %s", uploadID)
	data[1] = fmt.Sprintf("partNumberArr: %v", partNumberArr)
	data[2] = fmt.Sprintf("presignedURLs: %v", presignedURLs)
	common.Success(c, "Successlly generated presigned URL", data)
}

// POST /s3/upload/big/finish
func MultipartUploadFinish(c *gin.Context) {
	objectname := c.PostForm("objectname")
	uploadID := c.PostForm("uploadID")
	partNumberArr := c.PostFormArray("partNumberArr")
	partNumbers := make([]int, 0)
	for _, v := range partNumberArr {
		partNumber, _ := strconv.Atoi(v)
		partNumbers = append(partNumbers, partNumber)
	}
	sort.Ints(partNumbers)
	parts := make([]minio.CompletePart, 0)
	for _, v := range partNumbers {
		parts = append(parts, minio.CompletePart{
			PartNumber: v,
		})
	}
	uploadInfo, err := g.S3Client.CompleteMultipartUpload(g.S3Ctx, g.Config.Minio.BucketName, objectname, uploadID, parts, minio.PutObjectOptions{})
	if err != nil {
		common.Error(c, "生成预签名URL失败", err)
	} else {
		common.Success(c, "上传完成", fmt.Sprintf("uploadInfo: %v", uploadInfo))
	}
}

// POST /s3/upload/abort
func MultipartUploadAbort(c *gin.Context) {
	objectname := c.PostForm("objectname")
	uploadID := c.PostForm("uploadID")

	err := g.S3Client.AbortMultipartUpload(g.S3Ctx, g.Config.Minio.BucketName, objectname, uploadID)
	if err != nil {
		common.Error(c, "取消上传失败", err)
	} else {
		common.Success(c, "上传任务已取消")
	}
}

// GET /s3/upload/tasklist
func UploadTaskList(c *gin.Context) {
	result, err := g.S3Client.ListMultipartUploads(g.S3Ctx, g.Config.Minio.BucketName, "", "", "", "", 100)
	if err != nil {
		common.Error(c, "获取文件上传任务失败", err)
	} else {
		tasklist := fmt.Sprintf("uploads: %v",result.Uploads)
		common.Success(c, "获取上传任务成功", tasklist)
	}
}

// DownloadFile generates a pre-signed URL for downloading a file.
// GET /s3/download
func DownloadFile(c *gin.Context) {
	objectname := c.PostForm("objectname")

	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; file=\"%s\"", objectname))
	presignedURL, err := g.S3Client.PresignedGetObject(g.S3Ctx, g.Config.Minio.BucketName, objectname, DownloadExpiry, reqParams)
	if err != nil {
		common.Error(c, "生成预签名URL失败", err)
	}else{
		uri := fmt.Sprintf("downloadUri: %s",presignedURL.RequestURI())
		common.Success(c, "Successlly generated presigned URL", uri)
	}
}

// DELATE /s3/delate
func DelFile(c *gin.Context) {

}

// POST /s3/move
func MoveFile(c *gin.Context) {

}

// POST /s3/rename
func RenameFile(c *gin.Context) {

}

// POST /s3/mkdir
func Mkdir(c *gin.Context) {

}

// POST /s3/list
func FileList(c *gin.Context) {

}
