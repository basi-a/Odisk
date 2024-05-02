package controller

import (
	"fmt"
	"log"
	"strconv"

	"net/url"

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
	DefaultDownloadExpiry = time.Second * 24 * 60 * 60 * 7 // 7 days.
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
		return
	}
	presignedURL, err := g.S3core.Client.PresignedPutObject(g.S3Ctx, data.BucketName, data.ObjectName, UploadExpiry)
	if err != nil {
		common.Error(c, "生成预签名URL失败", err)
		return
	} else {
		common.Success(c, "Successlly generated presigned URL", map[string]string{"uploadUrl": presignedURL.String()})
	}

}

// POST /s3/upload/big/create
func MultipartUploadCreate(c *gin.Context) {

	type JsonData struct {
		BucketName    string `json:"bucketname"`
		ObjectName    string `json:"objectname"`
		MaxPartNumber int    `json:"maxPartNumber"` // min 1
	}
	data := JsonData{}

	if err := c.ShouldBindJSON(&data); err != nil {
		common.Error(c, "绑定失败", err)
		return
	}
	uploadID, err := g.S3core.NewMultipartUpload(g.S3Ctx, data.BucketName, data.ObjectName, minio.PutObjectOptions{})
	if err != nil {
		common.Error(c, "生成uploadID失败", err)
	}

	presignedURLs := make([]string, 0)

	// for  v := range data.PartNumberArr {
	for i := 1; i <= data.MaxPartNumber; i++ {
		// Get resources properly escaped and lined up before using them in http request.
		urlValues := make(url.Values)
		// Set part number.
		urlValues.Set("partNumber", strconv.Itoa(i))
		// Set upload id.
		urlValues.Set("uploadId", uploadID)
		presignedURL, err := g.S3core.Presign(g.S3Ctx, "PUT", data.BucketName, data.ObjectName, UploadExpiry, urlValues)
		if err != nil {
			common.Error(c, "生成预签名URL失败", err)
			return
		}
		presignedURLs = append(presignedURLs, presignedURL.String())
	}

	type Result struct {
		UploadID      string   `json:"uploadID"`
		PresignedURLs []string `json:"presignedURLs"`
	}
	result := Result{
		UploadID:      uploadID,
		PresignedURLs: presignedURLs,
	}
	common.Success(c, "Successlly generated presigned URL", result)
}

// POST /s3/upload/big/finish
func MultipartUploadFinish(c *gin.Context) {
	type JsonData struct {
		BucketName    string   `json:"bucketname"`
		ObjectName    string   `json:"objectname"`
		UploadID      string   `json:"uploadID"`
		MaxPartNumber int      `json:"maxPartNumber"`
		ETags         []string `json:"eTags"`
	}
	data := JsonData{}

	if err := c.ShouldBindJSON(&data); err != nil {
		common.Error(c, "绑定失败", err)
		return
	}

	parts := make([]minio.CompletePart, 0)
	for i := 1; i <= data.MaxPartNumber; i++ {
		parts = append(parts, minio.CompletePart{
			PartNumber: i,
			ETag:       data.ETags[i-1],
		})
	}
	// log.Println(len(parts))
	uploadInfo, err := g.S3core.CompleteMultipartUpload(g.S3Ctx, data.BucketName, data.ObjectName, data.UploadID, parts, minio.PutObjectOptions{})
	if err != nil {
		common.Error(c, "合并失败", err)
	} else {
		common.Success(c, "上传完成", fmt.Sprintf("uploadInfo: %v", uploadInfo))
	}
}

// POST /s3/download
func DownloadFile(c *gin.Context) {
	type JsonData struct {
		ObjectName     string        `json:"objectname"`
		BucketName     string        `json:"bucketname"`
		DownloadExpiry time.Duration `json:"downloadExpiry"`
	}
	data := JsonData{}
	if err := c.ShouldBindJSON(&data); err != nil {
		common.Error(c, "绑定失败", err)
		return
	}

	// Split the object 
	splitedObjectname := strings.Split(data.ObjectName, "/")
	
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", splitedObjectname[len(splitedObjectname)-1]))
	var presignedURL *url.URL
	var err error
	if data.DownloadExpiry != 0 {
		downloadExpiryDuration := time.Duration(data.DownloadExpiry) * time.Second
		presignedURL, err = g.S3core.PresignedGetObject(g.S3Ctx, data.BucketName, data.ObjectName, downloadExpiryDuration, reqParams)
		if err != nil {
			common.Error(c, "生成预签名URL失败", err)
			return
		}
	} else {
		presignedURL, err = g.S3core.PresignedGetObject(g.S3Ctx, data.BucketName, data.ObjectName, DefaultDownloadExpiry, reqParams)
		if err != nil {
			common.Error(c, "生成预签名URL失败", err)
			return
		}
	}

	common.Success(c, "Successlly generated presigned URL", map[string]string{"downloadUrl": presignedURL.String()})
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
		return
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
		return
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
		return
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
		return
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
				Size:         "",
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

// PUT /s3/upload/task/add
func TaskAdd(c *gin.Context) {
	type JsonData struct {
		BucketName string `json:"bucketname"`
		ObjectName string `json:"objectname"`
		UploadID   string `json:"uploadID"` // 小文件没这个
		SizeType   string `json:"sizetype"` // big or small
		Status     bool   `json:"status"`
	}
	data := JsonData{}

	if err := c.ShouldBindJSON(&data); err != nil {
		common.Error(c, "绑定失败", err)
		return
	}
	task := m.Task{
		BucketName: data.BucketName,
		ObjectName: data.ObjectName,
		UploadID:   data.UploadID,
		SizeType:   data.SizeType,
		Status:     data.Status,
	}
	if err := task.TaskAdd(); err != nil {
		common.Error(c, "记录任务失败", err)
	} else {
		common.Success(c, "记录任务成功", task.CreatedAt.Format("2006-01-02 15:04:05"))
	}

}

// PUT /s3/upload/task/done
func TaskDone(c *gin.Context) {
	type JsonData struct {
		BucketName string `json:"bucketname"`
		ObjectName string `json:"objectname"`
	}
	data := JsonData{}

	if err := c.ShouldBindJSON(&data); err != nil {
		common.Error(c, "绑定失败", err)
		return
	}
	task := m.Task{
		BucketName: data.BucketName,
		ObjectName: data.ObjectName,
	}
	if err := task.LocateTask(); err != nil {
		common.Error(c, "定位任务失败", err)
	}
	if err := task.TaskDone(); err != nil {
		common.Error(c, "任务状态标记失败", err)

	} else {
		common.Success(c, "任务状态标记成功", task.Status)
	}
}

// DELATE /s3/upload/task/del
func TaskDel(c *gin.Context) {
	type JsonData struct {
		BucketName string `json:"bucketname"`
		ObjectName string `json:"objectname"`
		UploadID   string `json:"uploadID"`
	}
	data := JsonData{}

	if err := c.ShouldBindJSON(&data); err != nil {
		common.Error(c, "绑定失败", err)
		return
	}
	if data.UploadID != "" {

		if err := g.S3core.AbortMultipartUpload(g.S3Ctx, data.BucketName, data.ObjectName, data.UploadID); err != nil {
			common.Error(c, "取消上传失败", err)
			return
		}
	}
	task := m.Task{
		BucketName: data.BucketName,
		ObjectName: data.ObjectName,
		UploadID:   data.UploadID,
	}

	if err := task.LocateTask(); err != nil {
		common.Error(c, "定位任务失败", err)
	}
	if err := task.TaskDel(); err != nil {
		common.Error(c, "任务删除/取消失败", err)

	} else {
		common.Success(c, "任务删除/取消成功", nil)
	}
}

// POST /s3/upload/task/list 这个列表要从数据库获取，minio不维护这个
func GetTaskList(c *gin.Context) {
	type JsonData struct {
		BucketName string `json:"bucketname"`
	}
	data := JsonData{}

	if err := c.ShouldBindJSON(&data); err != nil {
		common.Error(c, "绑定失败", err)
		return
	}
	bucketmap := m.Bucketmap{
		BucketName: data.BucketName,
		TaskList:   make([]m.Task, 0),
	}

	if bucketmap.BucketName != "" {
		// 先获取Bucketmap实例
		if err := bucketmap.GetMap(); err != nil {
			common.Error(c, "查找Bucketmap失败", err)
			return
		}
		if err := bucketmap.GetTaskList(); err != nil {
			common.Error(c, "获取列表失败", err)
			return
		}
	} else {
		if err := bucketmap.GetTaskListAll(); err != nil {
			common.Error(c, "获取列表失败", err)
			return
		}
	}
	common.Success(c, "获取列表成功", bucketmap.TaskList)
}

// DELATE /s3/bucketmapdel
func DeleteBucketMapWithTask(c *gin.Context) {
	type JsonData struct {
		BucketName string `json:"bucketname"`
	}
	data := JsonData{}
	if err := c.ShouldBindJSON(&data); err != nil {
		common.Error(c, "绑定失败", err)
		return
	}
	bucketmap := m.Bucketmap{
		BucketName: data.BucketName,
	}
	if err := bucketmap.GetMap(); err != nil {
		common.Error(c, "查找Bucketmap失败", err)
		return
	}

	if err := bucketmap.DeleteBucketMapWithTask(); err != nil {
		common.Error(c, "删除失败", err)
	} else {
		common.Success(c, "删除成功", nil)
	}
}
