package global

import (
	"context"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var S3core *minio.Core
var S3Ctx context.Context

func InitMinio() {
	S3Ctx = context.Background()
	RetryWithExponentialBackoff(UseMinio, "Minio Connection", 5)
}

func UseMinio() error {
	var err error
	S3core, err = minio.NewCore(Config.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(Config.Minio.AccessKeyID, Config.Minio.SecretAccessKey, ""),
		Secure: Config.Minio.UseSSL,
		Region: Config.Minio.Location,
	})
	if err != nil {
		return err
	}
	return nil
}

func MakeBucket(bucketName string) error {

	err := S3core.MakeBucket(S3Ctx, bucketName, minio.MakeBucketOptions{
		Region: Config.Minio.Location,
	})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, err := S3core.BucketExists(S3Ctx, bucketName)
		if err == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln("MiniO error:", err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}
	return nil
}
