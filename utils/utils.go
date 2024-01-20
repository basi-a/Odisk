package utils

import (

	"regexp"
)

// 检查邮箱格式
func IsValidEmail(email string) bool {
	// 邮箱正则表达式
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    return emailRegex.MatchString(email)
}

