package global

import (
	"log"
	"time"
)

type RetryFunc func() error

// RetryWithExponentialBackoff 是一个高阶函数，接受一个重试函数并执行重试逻辑
func RetryWithExponentialBackoff(fn RetryFunc, operationName string) {
	maxRetryCount := 5
	for retryCount := 0; retryCount < maxRetryCount; retryCount++ {
		err := fn()
		if err != nil {
			log.Printf("Retry failed for %s, error: %v. Waiting before retrying...\n", operationName, err)
			waitTime := time.Duration(retryCount*retryCount) * time.Second
			time.Sleep(waitTime)
		} else {
			log.Printf("Retry %s. Operation result: Success\n", operationName)
			return
		}
	}
	log.Fatalf("Failed after %d attempts for %s\n", maxRetryCount, operationName)
}