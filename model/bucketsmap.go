package model

import (
	g "odisk/global"

	"gorm.io/gorm"
)

type Bucketmap struct {
	gorm.Model

	UserID     uint   `gorm:"uniqueIndex"`//:idx_bucketmap_unique"`
	BucketName string `gorm:"uniqueIndex"`//:idx_bucketmap_unique;not null"`

	// 建立一对多关联关系，添加 onDelete: ReferentialAction.Cascade
	TaskList []Task `gorm:"foreignKey:BucketName;references:BucketName;constraint:OnDelete:CASCADE"`
}

type Task struct {
	gorm.Model
	BucketName string `json:"bucketname" gorm:"index:idx_task_bucketname;not null"`
	ObjectName string `json:"objectname"`
	UploadID   string `json:"uploadID"` // 小文件没这个
	SizeType   string `json:"sizetype"` // big or small
	Status     bool   `json:"status" gorm:"default:false"`
}

func AutoMigrateBucketmapAndTaskList() {
	if g.DB.Migrator().HasTable(&Bucketmap{}) && g.DB.Migrator().HasTable(&Task{}) {
		return
	}

	g.DB.AutoMigrate(&Bucketmap{})
	g.DB.AutoMigrate(&Task{})
}

func (bucketmap *Bucketmap) SaveMap() error {

	return g.DB.Create(&bucketmap).Error
}

func (bucketmap *Bucketmap) GetUserBucketName() error {

	return g.DB.Select("bucket_name").Where("user_id=?", bucketmap.UserID).First(&bucketmap).Error
}
func (bucketmap *Bucketmap) GetMap() error {

	return g.DB.Where("bucket_name =?", bucketmap.BucketName).First(&bucketmap).Error
}

func (bucketmap *Bucketmap) DeleteBucketMapWithTask() error {
	// 使用事务确保操作安全
	tx := g.DB.Begin()

	// 删除Bucketmap，由于设置了cascade，关联的Task会自动删除
	if err := tx.Delete(bucketmap).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	return tx.Commit().Error
}

func (task *Task)LocateTask() error {

	return g.DB.Where("bucket_name = ? AND object_name = ?", task.BucketName, task.ObjectName).First(&task).Error
}

func (task *Task) TaskDel() error {

	return g.DB.Delete(&task).Error
}
func (task *Task) TaskAdd() error {

	return g.DB.Create(&task).Error
}

func (task *Task) TaskDone() error {
	return g.DB.Model(&task).Update("status", true).Error
}


func (bucketmap *Bucketmap) GetTaskList() error {
	return g.DB.Where("bucket_name = ?", bucketmap.BucketName).Find(&bucketmap.TaskList).Error
}


func (bucketmap *Bucketmap) GetTaskListAll() error {

	return g.DB.Find(&bucketmap.TaskList).Error
}