package controller

import (
	"fmt"
	"log"

	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"odisk/common"
	g "odisk/global"
	m "odisk/model"
	u "odisk/utils"

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
	type JsonData struct {
		ObjectName string `json:"objectname"`
		BucketName string `json:"bucketname"`
	}
	data := JsonData{}
	if err := c.ShouldBindJSON(&data); err != nil {
		common.Error(c, "绑定失败", err)
	}
	presignedURL, err := g.S3core.Client.PresignedPutObject(g.S3Ctx, data.BucketName, data.ObjectName, UploadExpiry)
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

	uploadID, err := g.S3core.NewMultipartUpload(g.S3Ctx, bucketname, objectname, minio.PutObjectOptions{})
	if err != nil {
		common.Error(c, "生成uploadID失败", err)
	}

	presignedURLs := make([]string, 0)
	for _, v := range partNumberArr {

		presignedURL, err := g.S3core.Presign(g.S3Ctx, "PUT", bucketname, objectname, UploadExpiry, url.Values{
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
	uploadInfo, err := g.S3core.CompleteMultipartUpload(g.S3Ctx, bucketname, objectname, uploadID, parts, minio.PutObjectOptions{})
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
	err := g.S3core.AbortMultipartUpload(g.S3Ctx, bucketname, objectname, uploadID)
	if err != nil {
		common.Error(c, "取消上传失败", err)
	} else {
		common.Success(c, "上传任务已取消", nil)
	}
}

// GET /s3/upload/tasklist
func UploadTaskList(c *gin.Context) {
	bucketname := c.PostForm("bucketname")
	prefix := c.DefaultPostForm("prefix", "")
	result, err := g.S3core.ListMultipartUploads(g.S3Ctx, bucketname, prefix, "", "", "", 100)
	if err != nil {
		common.Error(c, "获取文件上传任务失败", err)
	} else {
		tasklist := fmt.Sprintf("uploads: %v", result.Uploads)
		common.Success(c, "获取上传任务成功", tasklist)
	}
}

// POST /s3/download
func DownloadFile(c *gin.Context) {
	type JsonData struct {
		ObjectName string `json:"objectname"`
		BucketName string `json:"bucketname"`
	}
	data := JsonData{}
	if err := c.ShouldBindJSON(&data); err != nil {
		common.Error(c, "绑定失败", err)
	}
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", data.ObjectName))

	presignedURL, err := g.S3core.PresignedGetObject(g.S3Ctx, data.BucketName, data.ObjectName, DownloadExpiry, reqParams)
	if err != nil {
		common.Error(c, "生成预签名URL失败", err)
		return
	}

	common.Success(c, "Successlly generated presigned URL", map[string]string{"downloadUri": presignedURL.String()})
}

// DELETE /s3/delate
func DeleteFile(c *gin.Context) {

	type JsonData struct {
		ObjectName string `json:"objectname"`
		BucketName string `json:"bucketname"`
	}
	data := JsonData{}
	if err := c.ShouldBindJSON(&data); err != nil {
		common.Error(c, "绑定失败", err)
	}
	err := g.S3core.RemoveObject(g.S3Ctx, data.BucketName, data.ObjectName, minio.RemoveObjectOptions{})
	if err != nil {
		common.Error(c, "删除文件失败", err)
		return
	}

	common.Success(c, "文件删除成功", nil)
}

// POST /s3/mv
func MoveFile(c *gin.Context) {
	type JsonData struct {
		SrcBucketName  string `json:"srcbucketname"`
		SrcObjectName  string `json:"srcobjectname"`
		DestObjectName string `json:"destobjectName"`
	}

	data := JsonData{}
	if err := c.ShouldBindJSON(&data); err != nil {
		common.Error(c, "绑定失败", err)
	}
	_, err := g.S3core.Client.StatObject(g.S3Ctx, data.SrcBucketName, data.SrcObjectName, minio.StatObjectOptions{})
	if err != nil {
		common.Error(c, "源文件不存在或无法访问", err)
		return
	}
	// Split the destination object name into segments
	destSegments := strings.Split(data.DestObjectName, "/")

	// Check if the target directory exists
	parentDir := strings.Join(destSegments[:len(destSegments)-1], "/")
	if !strings.HasSuffix(parentDir, "/") {
		parentDir += "/" // Ensure each directory segment ends with a slash
	}

	_, err = g.S3core.Client.StatObject(g.S3Ctx, data.SrcBucketName, parentDir, minio.StatObjectOptions{})
	if err != nil {
		// Target directory doesn't exist, return an error message
		common.Error(c, "目标目录:"+parentDir+" , 不存在, 请先创建", err)
		return
	}
	copyDest := minio.CopyDestOptions{
		Bucket: data.SrcBucketName,
		Object: data.DestObjectName,
	}
	copySrc := minio.CopySrcOptions{
		Bucket: data.SrcBucketName,
		Object: data.SrcObjectName,
	}

	_, err = g.S3core.Client.CopyObject(g.S3Ctx, copyDest, copySrc)
	if err != nil {
		common.Error(c, "复制文件失败", err)
		return
	}

	err = g.S3core.Client.RemoveObject(g.S3Ctx, data.SrcBucketName, data.SrcObjectName, minio.RemoveObjectOptions{})
	if err != nil {
		common.Error(c, "删除源文件失败（复制成功，但源文件未删除）", err)
		return
	}

	common.Success(c, "文件移动(重命名)成功", nil)
}

// POST /s3/mkdir
func Mkdir(c *gin.Context) {

	type JsonData struct {
		BucketName string `json:"bucketname"`
		DirName    string `json:"dirname"`
	}
	data := JsonData{}
	if err := c.ShouldBindJSON(&data); err != nil {
		common.Error(c, "绑定失败", err)
	}
	// Ensure the directory name ends with a slash
	if !strings.HasSuffix(data.DirName, "/") {
		data.DirName += "/"
	}

	// Split the directory name into segments
	segments := strings.Split(data.DirName, "/")

	/// Check and create parent directories and the final directory
	for i := 1; i < len(segments); i++ {
		parentDir := strings.Join(segments[:i], "/")
		if !strings.HasSuffix(parentDir, "/") {
			parentDir += "/" // Ensure each directory segment ends with a slash
		}
		log.Println(parentDir)
		_, err := g.S3core.StatObject(g.S3Ctx, data.BucketName, parentDir, minio.StatObjectOptions{})
		if err != nil {
			// Parent directory doesn't exist, create it
			_, err := g.S3core.PutObject(g.S3Ctx, data.BucketName, parentDir, strings.NewReader(""), 0, "", "", minio.PutObjectOptions{})
			if err != nil {
				common.Error(c, "目录创建失败", err)
				return
			}
		}
	}

	common.Success(c, "目录创建成功", data.DirName)
}

// post /s3/list
func FileList(c *gin.Context) {
	type JsonData struct {
		BucketName string `json:"bucketname"`
		Prefix     string `json:"prefix"`
	}
	data := JsonData{}

	if err := c.ShouldBindJSON(&data); err != nil {
		common.Error(c, "绑定失败", err)
	}

	ch := g.S3core.Client.ListObjects(g.S3Ctx, data.BucketName, minio.ListObjectsOptions{
		Prefix:    data.Prefix,
		Recursive: false,
	})

	fileInfos := []m.FileInfo{}

	// 使用for-range遍历通道
	for v := range ch {
		isdir := false
		if v.Size == 0 && strings.HasSuffix(v.Key, "/") {
			isdir = !isdir
		}
		if isdir {
			fileInfo := m.FileInfo{
				Key:          v.Key,
				LastModified: "",
				Size:         u.FormatFileSize(v.Size),
				ContentType:  "directory",
				IsDir:        isdir,
			}
			fileInfos = append(fileInfos, fileInfo)
		} else {
			fileInfo := m.FileInfo{
				Key:          v.Key,
				LastModified: v.LastModified.Format("2006-01-02 15:04:05"),
				Size:         u.FormatFileSize(v.Size),
				ContentType:  u.GuessContentTypeFromExtension(v.Key),
				IsDir:        isdir,
			}
			fileInfos = append(fileInfos, fileInfo)
		}

	}

	// 返回结果
	common.Success(c, "获取文件列表成功", fileInfos)
}
