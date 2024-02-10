package utils

import (
	"encoding/base64"
	"github.com/wneessen/go-mail"
	"regexp"
	"html/template"
	g "odisk/global"
)

// 检查邮箱格式
func IsValidEmail(email string) bool {
	// 邮箱正则表达式
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    return emailRegex.MatchString(email)
}

// 发送邮件
func SendEmail(email string, subject string, t *template.Template, data interface{}) error{	
	msg := mail.NewMsg()
	if err := msg.From(g.Config.Server.Mail.SerderMail); err != nil {
		return err
	}
	if err := msg.To(email); err != nil {
		return err
	}
	msg.Subject(subject)
	msg.SetBodyHTMLTemplate(t, data)

	client, err := mail.NewClient(g.Config.Server.Mail.SmtpServer, 
		mail.WithPort(g.Config.Server.Mail.Port), 
		mail.WithSMTPAuth(mail.SMTPAuthPlain), 
		mail.WithUsername(g.Config.Server.Mail.UserName), 
		mail.WithPassword(g.Config.Server.Mail.Password))
	if err != nil {
		return err
	}
	client.DialAndSend(msg)
	return nil
}
// base64解编码
func DecodeRawData(encodedRawData string) (string, error){
	decodedstr, err := base64.StdEncoding.DecodeString(encodedRawData)
	if err != nil {
		return "", err
	}
	return string(decodedstr), nil
}