package model

import (

	g "odisk/global"

	"gorm.io/gorm"
)

type Bucketmap struct {
	gorm.Model

	UserID     uint   `gorm:"uniqueIndex"`
	BucketName string `gorm:"uniqueIndex"`

	// 建立一对多关联关系
	TaskList []Task `gorm:"foreignKey:BucketName;references:BucketName"`
}

type Task struct {
	gorm.Model
	BucketName string `json:"bucketname" gorm:"index:idx_task_bucketname;not null"`
	ObjectName string `json:"objectname" gorm:"not null"`
	FileName   string `json:"filename"`
	UploadID   string `json:"uploadID"` // 小文件没这个
	Size       uint   `json:"size"`
	Status     string   `json:"status" gorm:"default:uploading;not null"` // uploading done removed
}

func AutoMigrateBucketmapAndTaskList() {

	if !g.DB.Migrator().HasTable(&Bucketmap{}) {
		g.DB.AutoMigrate(&Bucketmap{})
	}
	if !g.DB.Migrator().HasTable(&Task{}) {
		g.DB.AutoMigrate(&Task{})
	}
}

func (bucketmap *Bucketmap) SaveMap() error {
	return g.DB.Create(&bucketmap).Error
}

func (bucketmap *Bucketmap) GetUserBucketName() error {

	return g.DB.Select("bucket_name").Where("user_id=?", bucketmap.UserID).First(&bucketmap).Error
}
func (bucketmap *Bucketmap) GetMap() error {

	return g.DB.Where("user_id=?", bucketmap.UserID).First(&bucketmap).Error
}
func (bucketmap *Bucketmap) GetMapByBucketName() error {

	return g.DB.Where("bucket_name=?", bucketmap.BucketName).First(&bucketmap).Error
}
func GetMapList() ([]Bucketmap, error)  {
	var mapList []Bucketmap
	if err := g.DB.Find(&mapList).Error; err != nil {
		return []Bucketmap{}, err
	}
	return mapList, nil
}

func (bucketmap *Bucketmap) UpdateMap(NewBucketName string) error {
	// 开始事务
	tx := g.DB.Begin()

	// 更新关联的Task的BucketName
	if err := tx.Model(&Task{}).Where("bucket_name = ?", bucketmap.BucketName).Update("bucket_name", NewBucketName).Error; err != nil {
		tx.Rollback() // 出错则回滚
		return err
	}

	// 更新Bucketmap的bucket_name
	if err := tx.Model(&bucketmap).Where("user_id = ?", bucketmap.UserID).Update("bucket_name", NewBucketName).Error; err != nil {
		tx.Rollback() // 更新Bucketmap出错也回滚
		return err
	}

	// 提交事务
	return tx.Commit().Error
}

func (bucketmap *Bucketmap) DeleteBucketMapWithTask() error {
	// 使用事务确保操作安全
	tx := g.DB.Begin()
	// 手动删除关联的Task
    if err := tx.Where("bucket_name = ?", bucketmap.BucketName).Delete(&Task{}).Error; err != nil {
        tx.Rollback()
        return err
    }
	// 删除Bucketmap，由于设置了cascade，关联的Task会自动删除
	if err := tx.Delete(bucketmap).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	return tx.Commit().Error
}

func (task *Task) LocateTask(id uint) error {
	return g.DB.Where("id = ?", id).First(&task).Error
}

func (task *Task) TaskDel(id uint) error {
	return g.DB.Where("id = ?", id).Delete(&task).Error
}
func (task *Task) TaskAdd() error {

	return g.DB.Create(&task).Error
}

func (task *Task) TaskDone(id uint) error {
	return g.DB.Model(&task).Where("id = ?", id).Update("status", "done").Error
}
func (task *Task) TaskAbort(id uint) error {
	return g.DB.Model(&task).Where("id = ?", id).Update("status", "removed").Error
}


func (bucketmap *Bucketmap) GetTaskList() error {
	return g.DB.Where("bucket_name = ?", bucketmap.BucketName).Find(&bucketmap.TaskList).Error
}

func (bucketmap *Bucketmap) GetTaskListAll() error {

	return g.DB.Find(&bucketmap.TaskList).Error
}
