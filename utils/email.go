package utils

import (
	"encoding/base64"
	"regexp"
)

// 检查邮箱格式
func IsValidEmail(email string) bool {
	// 邮箱正则表达式
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    return emailRegex.MatchString(email)
}

// base64解编码
func DecodeRawData(encodedRawData string) (string, error){
	decodedstr, err := base64.StdEncoding.DecodeString(encodedRawData)
	if err != nil {
		return "", err
	}
	return string(decodedstr), nil
}