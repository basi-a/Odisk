package global

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"math"
	"net/http"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var S3core *minio.Core
var S3Ctx context.Context

func InitMinio() {
	S3Ctx = context.Background()
	RetryWithExponentialBackoff(UseMinio, "Minio Connection", 5)
}

func MaxBucketSize() int {
	max := Config.Minio.BucketMaxSize* int(math.Pow(2,30))
	// log.Println(max)
	return max
}

func UseMinio() error {
	var err error
	// 假设您有一个或多个CA证书文件（例如.crt文件）
	certPool := x509.NewCertPool()

	// 从文件中读取证书
	caCertPath := Config.Server.Ssl.Cert // 更改为实际的证书路径
	caCertBytes, err := os.ReadFile(caCertPath)
	if err != nil {
		log.Fatalf("could not read CA certificate: %v", err)
	}

	// 尝试将PEM编码的证书添加到证书池
	if ok := certPool.AppendCertsFromPEM(caCertBytes); !ok {
		log.Fatalf("failed to append certificate to pool")
	}
	httpTransport := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: certPool, // 设置证书池
		},
	}
	S3core, err = minio.NewCore(Config.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(Config.Minio.AccessKeyID, Config.Minio.SecretAccessKey, ""),
		Secure: Config.Minio.UseSSL,
		Transport: httpTransport,
		Region: Config.Minio.Location,
	})
	if err != nil {
		return err
	}
	return nil
}

func MakeBucket(bucketName string) error {

	err := S3core.Client.MakeBucket(S3Ctx, bucketName, minio.MakeBucketOptions{
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

func GetCurrentSize(bucketName string) int{
	ch := S3core.Client.ListObjects(S3Ctx, bucketName, minio.ListObjectsOptions{
		Recursive: false,
	})
	var currentBucketSize int
	for v := range ch {
		currentBucketSize += int(v.Size)
	}
	return currentBucketSize
}