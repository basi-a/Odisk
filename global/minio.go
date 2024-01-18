package global

import (
	"log"
	"time"

	"github.com/minio/minio-go"
)
var MinioClient *minio.Client

func InitMinio()  {
	endpoint 		:= Config.Minio.Endpoint
	accessKeyId 	:= Config.Minio.AccessKeyID
	secretAccessKey := Config.Minio.SecretAccessKey
	usessl 			:= Config.Minio.UseSSL
	bucketName 		:= Config.Minio.BucketName
	location		:= Config.Minio.Location

	maxRetryCount := 5
	var err error
	for retryCount := 0; retryCount < maxRetryCount; retryCount++{
		MinioClient, err = minio.New(endpoint, accessKeyId, secretAccessKey, usessl)	
		if err == nil {
			break
		}else {
			log.Println("mariadb error:",err)
		}
		time.Sleep(time.Second*20)
	}
	
	err = MinioClient.MakeBucket(bucketName, location)
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, err := MinioClient.BucketExists(bucketName)
		if err == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln("MiniO error:",err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}
}