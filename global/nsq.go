package global

import (
	"encoding/json"

	"html/template"

	"strings"

	"github.com/nsqio/go-nsq"
	"github.com/wneessen/go-mail"
)

type EmailData struct {
	Email   string             `json:"email"`
	Subject string             `json:"subject"`
	T       *template.Template `json:"t"`
	Data    interface{}        `json:"data"`
}

var Producer *nsq.Producer
var Consumer *nsq.Consumer

func InitNsq() {

	RetryWithExponentialBackoff(
		func() error {
			return CreateAndStartConsumer("email", "user-auth")
		}, "Create and start Nsq Consumer")

	RetryWithExponentialBackoff(CreateProducer, "Create Nsq Producer")
}

func CreateProducer() error {
	var err error
	nsqdAddr := strings.Join(Config.Nsq.Nsqd, ";")
	Producer, err = nsq.NewProducer(nsqdAddr, nsq.NewConfig())
	if err != nil {
		return err
	}
	return nil
}

func CreateAndStartConsumer(topic, channel string) error {
	var err error

	Consumer, err = nsq.NewConsumer(topic, channel, nsq.NewConfig())
	if err != nil {
		return err
	}

	Consumer.AddConcurrentHandlers(nsq.HandlerFunc(msgHandler), 5)

	if err = Consumer.ConnectToNSQLookupds(Config.Nsq.Nsqlookupd); err != nil {
		return err
	}
	return nil
}

func NsqPublish(topic string, messgageBody interface{}) error {
	body, err := json.Marshal(messgageBody)
	if err != nil {
		return err
	}
	return Producer.Publish(topic, body)
}

func msgHandler(message *nsq.Message) error {
	var data interface{}
	err := json.Unmarshal(message.Body, &data)
	if err != nil {
		return err
	}

	switch ty := data.(type) {
	case EmailData:
		if err := sendEmail(ty.Email, ty.Subject, ty.T, ty.Data); err != nil {
			return err
		}
	default:
		return &UnknowTypeError{ty}
	}

	return nil
}

/*
	SendEmail 发送一封邮件给指定的收件人。

参数：
email: 收件人的邮箱地址。
subject: 邮件的主题。
t: 用于生成邮件正文的模板。
data: 模板所需的数据。
返回值：
error: 如果在发送邮件过程中发生错误，则返回该错误；否则返回nil。
*/
func sendEmail(email string, subject string, t *template.Template, data interface{}) error {
	msg := mail.NewMsg()
	if err := msg.From(Config.Server.Mail.SerderMail); err != nil {
		return err
	}
	if err := msg.To(email); err != nil {
		return err
	}
	// 设置邮件的主题
	msg.Subject(subject)
	// 使用HTML模板和数据设置邮件正文
	msg.SetBodyHTMLTemplate(t, data)

	client, err := mail.NewClient(Config.Server.Mail.SmtpServer,
		mail.WithPort(Config.Server.Mail.Port),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(Config.Server.Mail.UserName),
		mail.WithPassword(Config.Server.Mail.Password))
	if err != nil {
		return err
	}
	client.DialAndSend(msg)
	return nil
}
