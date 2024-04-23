package controller

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
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
	bucketname := c.PostForm("bucketname")
	presignedURL, err := g.S3Client.PresignedPutObject(g.S3Ctx, bucketname, objectname, UploadExpiry)
	if err != nil {
		common.Error(c, "生成预签名URL失败", err)
		return
	} else {

		common.Success(c, "Successlly generated presigned URL", map[string]string{"uploadUri": presignedURL.String()})
	}

}

// POST /s3/upload/big/create
func MultipartUploadCreate(c *gin.Context) {
	objectname := c.PostForm("objectname")
	bucketname := c.PostForm("bucketname")
	partNumberArr := c.PostFormArray("partNumberArr")

	uploadID, err := g.S3Client.NewMultipartUpload(g.S3Ctx, bucketname, objectname, minio.PutObjectOptions{})
	if err != nil {
		common.Error(c, "生成uploadID失败", err)
	}

	presignedURLs := make([]string, 0)
	for _, v := range partNumberArr {

		presignedURL, err := g.S3Client.Presign(g.S3Ctx, "PUT", bucketname, objectname, UploadExpiry, url.Values{
			"uploadID":   []string{uploadID},
			"partNumber": []string{v},
		})
		if err != nil {
			common.Error(c, "生成预签名URL失败", err)
			return
		}
		presignedURLs = append(presignedURLs, presignedURL.String())
	}
	data := make([]string, 3)
	data[0] = fmt.Sprintf("uploadID: %s", uploadID)
	data[1] = fmt.Sprintf("partNumberArr: %v", partNumberArr)
	data[2] = fmt.Sprintf("presignedURLs: %v", presignedURLs)
	common.Success(c, "Successlly generated presigned URL", data)
}

// POST /s3/upload/big/finish
func MultipartUploadFinish(c *gin.Context) {
	objectname := c.PostForm("objectname")
	uploadID := c.PostForm("uploadID")
	bucketname := c.PostForm("bucketname")
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
	uploadInfo, err := g.S3Client.CompleteMultipartUpload(g.S3Ctx, bucketname, objectname, uploadID, parts, minio.PutObjectOptions{})
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
	bucketname := c.PostForm("bucketname")
	err := g.S3Client.AbortMultipartUpload(g.S3Ctx, bucketname, objectname, uploadID)
	if err != nil {
		common.Error(c, "取消上传失败", err)
	} else {
		common.Success(c, "上传任务已取消")
	}
}

// GET /s3/upload/tasklist
func UploadTaskList(c *gin.Context) {
	bucketname := c.PostForm("bucketname")
	prefix := c.DefaultPostForm("prefix", "")
	result, err := g.S3Client.ListMultipartUploads(g.S3Ctx, bucketname, prefix, "", "", "", 100)
	if err != nil {
		common.Error(c, "获取文件上传任务失败", err)
	} else {
		tasklist := fmt.Sprintf("uploads: %v", result.Uploads)
		common.Success(c, "获取上传任务成功", tasklist)
	}
}

// POST /s3/download
func DownloadFile(c *gin.Context) {
	objectName := c.PostForm("objectname")
	bucketName := c.PostForm("bucketname")
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", objectName))

	presignedURL, err := g.S3Client.PresignedGetObject(g.S3Ctx, bucketName, objectName, DownloadExpiry, reqParams)
	if err != nil {
		common.Error(c, "生成预签名URL失败", err)
		return
	}

	common.Success(c, "Successlly generated presigned URL", map[string]string{"downloadUri": presignedURL.String()})
}

// DELETE /s3/delate
func DeleteFile(c *gin.Context) {
	objectName := c.PostForm("objectname")
	bucketName := c.PostForm("bucketname")

	err := g.S3Client.RemoveObject(g.S3Ctx, bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		common.Error(c, "删除文件失败", err)
		return
	}

	common.Success(c, "文件删除成功")
}

// POST /s3/mv
func MoveFile(c *gin.Context) {
	currentObjectName := c.PostForm("current_objectname")
	currentBucketName := c.PostForm("current_bucketname")
	newObjectName := c.PostForm("new_objectname")

	_, err := g.S3Client.StatObject(g.S3Ctx, currentBucketName, currentObjectName, minio.StatObjectOptions{})
	if err != nil {
		common.Error(c, "源文件不存在或无法访问", err)
		return
	}

	copySrc := minio.CopySrcOptions{
		Bucket: currentBucketName,
		Object: currentObjectName,
	}

	_, err = g.S3Client.CopyObject(g.S3Ctx, currentBucketName, currentObjectName, currentBucketName, newObjectName, nil, copySrc, minio.PutObjectOptions{})
	if err != nil {
		common.Error(c, "复制文件失败", err)
		return
	}

	err = g.S3Client.RemoveObject(g.S3Ctx, currentBucketName, currentObjectName, minio.RemoveObjectOptions{})
	if err != nil {
		common.Error(c, "删除源文件失败（复制成功，但源文件未删除）", err)
		return
	}

	common.Success(c, "文件移动(重命名)成功")
}

// POST /s3/mkdir
func Mkdir(c *gin.Context) {
	bucketName := c.PostForm("bucketname")
	directoryName := c.PostForm("dirname") + "/"

	// Ensure the directory name ends with a slash
	if !strings.HasSuffix(directoryName, "/") {
		directoryName += "/"
	}

	// Create an empty object with the given directory name
	_, err := g.S3Client.PutObject(g.S3Ctx, bucketName, directoryName, strings.NewReader(""), 0, "", "", minio.PutObjectOptions{})

	if err != nil {
		common.Error(c, "创建目录失败", err)
		return
	}

	common.Success(c, "目录创建成功", directoryName)
}

// POST /s3/list
func FileList(c *gin.Context) {
	bucketname := c.PostForm("bucketname")
	prefix := c.DefaultPostForm("prefix", "")

	// List objects with the specified prefix (virtual directory path)
	objects, err := g.S3Client.ListObjects(bucketname, prefix, "", "", 100)
	if err != nil {
		common.Error(c, "获取文件列表失败", err)
		return
	}

	// Prepare the response data
	fileList := make([]string, len(objects.Contents))
	for i, object := range objects.Contents {
		fileList[i] = object.Key
	}

	common.Success(c, "获取文件列表成功", fileList)
}
