package controller
import (
	"fmt"


	"strconv"

	"odisk/common"
	g "odisk/global"
	m "odisk/model"


	"github.com/gin-gonic/gin"

)
// PUT /v1/s3/upload/task/add
func TaskAdd(c *gin.Context) {
	type JsonData struct {
		BucketName string `json:"bucketname"`
		ObjectName string `json:"objectname"`
		FileName   string `json:"filename"`
		UploadID   string `json:"uploadID"` // 小文件没这个
		Size       uint   `json:"size"`
	}
	data := JsonData{}

	if err := c.ShouldBindJSON(&data); err != nil {
		common.Error(c, "绑定失败", err)
		return
	}
	task := m.Task{
		BucketName: data.BucketName,
		ObjectName: data.ObjectName,
		FileName:   data.FileName,
		UploadID:   data.UploadID,
		Size:       data.Size,
	}
	if err := task.TaskAdd(); err != nil {
		common.Error(c, "记录任务失败", err)
	} else {
		common.Success(c, "记录任务成功", map[string]uint{
			"taskID": task.ID,
		})
	}
}

// PUT /v1/s3/upload/task/percent/update
func UpdateTaskPercent(c *gin.Context) {
	type JsonData struct {
		TaskID  int `json:"taskID"`
		Percent int `json:"percent"`
	}
	data := JsonData{}

	if err := c.ShouldBindJSON(&data); err != nil {
		common.Error(c, "绑定失败", err)
		return
	}
	key := fmt.Sprintf("TaskPercent: %d", data.TaskID)
	SaveSession(c, key, data.Percent)
	common.Success(c, "更新成功", nil)
}

// GET /v1/s3/upload/task/percent/:taskID
func GetTaskPercent(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("taskID"))
	if err != nil {
		common.Error(c, "进度获取失败", err)
		return
	}
	// log.Println(taskID)
	key := fmt.Sprintf("TaskPercent: %d", taskID)
	value := ReadSession(c, key)
	// log.Println(value)
	if percent, ok := value.(int); ok {
		common.Success(c, "进度获取成功", map[string]int{
			"percent": percent,
		})
	} else {
		common.Error(c, "获取失败", nil)
	}
}

// PUT /v1/s3/upload/task/done
func TaskDone(c *gin.Context) {
	type JsonData struct {
		TaskID int `json:"taskID"`
	}
	data := JsonData{}

	if err := c.ShouldBindJSON(&data); err != nil {
		common.Error(c, "绑定失败", err)
		return
	}
	task := m.Task{}

	if err := task.TaskDone(uint(data.TaskID)); err != nil {
		common.Error(c, "任务状态标记失败", err)

	} else {
		common.Success(c, "任务状态标记成功", task.Status)
		key := fmt.Sprintf("TaskPercent: %d", task.ID)
		DelSession(c, key)
	}
}

// POST /v1/s3/upload/task/abort
func TaskAbort(c *gin.Context) {
	type JsonData struct {
		BucketName string `json:"bucketname"`
		ObjectName string `json:"objectname"`
		UploadID   string `json:"uploadID"`
		TaskID     int    `json:"taskID"`
	}
	data := JsonData{}

	if err := c.ShouldBindJSON(&data); err != nil {
		common.Error(c, "绑定失败", err)
		return
	}
	// log.Println(data)
	if data.UploadID != "----" {

		if err := g.S3core.AbortMultipartUpload(g.S3Ctx, data.BucketName, data.ObjectName, data.UploadID); err != nil {
			common.Error(c, "取消上传失败", err)
			return
		}
	}
	task := m.Task{}

	if err := task.TaskAbort(uint(data.TaskID)); err != nil {
		common.Error(c, "任务取消失败", err)
	} else {
		common.Success(c, "任务取消成功", nil)
	}
}

// POST /v1/s3/upload/task/del
func TaskDel(c *gin.Context) {
	type JsonData struct {
		TaskID int `json:"taskID"`
	}
	data := JsonData{}

	if err := c.ShouldBindJSON(&data); err != nil {
		common.Error(c, "绑定失败", err)
		return
	}
	task := m.Task{}
	if err := task.LocateTask(uint(data.TaskID)); err != nil {
		common.Error(c, "任务定位失败", err)
		return
	}

	if task.Status == "uploading" {
		common.Error(c, "上传中任务不能删除, 请先取消", nil)
		return
	}
	if err := task.TaskDel(uint(data.TaskID)); err != nil {
		common.Error(c, "任务取消失败", err)
		return
	} else {
		common.Success(c, "任务取消成功", nil)
	}
}

// POST /v1/s3/upload/task/list 这个列表要从数据库获取，minio不维护这个
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
		if err := bucketmap.GetMapByBucketName(); err != nil {
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

// POST /v1/s3/bucketmap/del
func DeleteBucketMapWithTask(c *gin.Context) {
	type JsonData struct {
		UserID     int    `json:"userID"`
		BucketName string `json:"bucketname"`
	}
	data := JsonData{}
	if err := c.ShouldBindJSON(&data); err != nil {
		common.Error(c, "绑定失败", err)
		return
	}
	bucketmap := m.Bucketmap{
		UserID:     uint(data.UserID),
		BucketName: data.BucketName,
	}
	if err := bucketmap.GetMap(); err != nil {
		common.Error(c, "查找Bucketmap失败", err)
		return
	}
	if err := g.DeactivateBucket(bucketmap.BucketName); err != nil {
		common.Error(c, "停用桶失败", err)
		return
	}
	if err := bucketmap.DeleteBucketMapWithTask(); err != nil {
		common.Error(c, "删除失败", err)
	} else {
		common.Success(c, "删除成功", nil)
	}
}



// GET /v1/s3/bucketmap/list
func GetMapList(c *gin.Context) {
	List, err := m.GetMapList()
	if err != nil {
		common.Error(c, "获取失败", err)
		return
	}
	common.Success(c, "获取成功", List)
}

// POST /v1/s3/bucketmap/update
func UpdateBucketmap(c *gin.Context) {
	type JsonData struct {
		UserID        int    `json:"userID"`
		BucketName    string `json:"bucketname"`
		NewBucketName string `json:"newBucketname"`
	}
	data := JsonData{}
	if err := c.ShouldBindJSON(&data); err != nil {
		common.Error(c, "绑定失败", err)
		return
	}
	bucketmap := m.Bucketmap{
		UserID:     uint(data.UserID),
		BucketName: data.BucketName,
	}
	if err := bucketmap.UpdateMap(data.NewBucketName); err != nil {
		common.Error(c, "更新失败", err)
		return
	}
	common.Success(c, "更新成功", nil)
}
