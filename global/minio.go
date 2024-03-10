package global

import (
	"context"
	"log"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)
var S3Client *minio.Core
var S3Ctx context.Context
func InitMinio()  {

	S3Ctx = context.Background()
	maxRetryCount := 5
	var err error
	for retryCount := 0; retryCount < maxRetryCount; retryCount++{
		S3Client, err = minio.NewCore(Config.Minio.Endpoint, &minio.Options{
			Creds: credentials.NewStaticV4(Config.Minio.AccessKeyID, Config.Minio.SecretAccessKey, ""),
			Secure: Config.Minio.UseSSL,
			Region: Config.Minio.Location,
		})
		if err == nil {
			break
		}else {
			log.Println("minio error:",err)
		}
		time.Sleep(time.Second*20)
	}
	

}

func MakeBucket(bucketName string)  error{

	err := S3Client.MakeBucket(S3Ctx, bucketName, minio.MakeBucketOptions{
		Region: Config.Minio.Location,
	})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, err := S3Client.BucketExists(S3Ctx, bucketName)
		if err == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln("MiniO error:",err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}
	return nil
}