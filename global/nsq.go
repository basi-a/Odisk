package global

import (
	"encoding/json"
	"fmt"
	"log"
	// "math/rand"

	"html/template"

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
var Consumers []*nsq.Consumer

func InitNsq() {
	defer log.Println("nsq producer and consumer initialization completed.")
	log.Println("Initialize nsq's producer and consumers....")
	RetryWithExponentialBackoff(CreateNsqProducer, "Create Nsq Producer", 5)
	RetryWithExponentialBackoff(CreateAndStartNsqConsumer, "Create and start Nsq Consumer", 5)
}

// 一个服务端程序，对应一个localhost的nsqd，且写入本地文件，然后集体使用nsqlookupd，均衡负载消费
func CreateNsqProducer() error {
	var err error

	if Producer, err = nsq.NewProducer(Config.Nsq.Nsqd+":"+Config.Nsq.Port.Nsqd.TCP, nsq.NewConfig()); err != nil {
		log.Println("Error creating producer:", err)
		return err
	}

	if err = Producer.Ping(); err != nil {
		log.Println("Error pinging nsqd:", err)
		Producer.Stop()
		return err
	}

	for topic := range Config.Nsq.Topics {
		if err = NsqPublish(topic, "topic created"); err != nil {
			log.Println("Error publishing message:", err)
			return err
		}
	}

	return nil
}

func CreateAndStartNsqConsumer() error {
	nsqlookupdAddrsWithPort := make([]string, 0)

	for _, v := range Config.Nsq.Nsqlookupd {
		nsqlookupdAddrsWithPort = append(nsqlookupdAddrsWithPort, v+":"+Config.Nsq.Port.Nsqlookupd.HTTP)
	}
	consumers := make([]*nsq.Consumer, 0)
	for topic, channals := range Config.Nsq.Topics {

		for _, channal := range channals {
			consumer, err := nsq.NewConsumer(topic, channal, nsq.NewConfig())
			if err != nil {
				return err
			}

			consumer.AddConcurrentHandlers(nsq.HandlerFunc(msgHandler), 5)

			if err = consumer.ConnectToNSQLookupds(nsqlookupdAddrsWithPort); err != nil {
				return err
			}

			consumers = append(consumers, consumer)
		}
	}
	Consumers = consumers
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
	case string:
		log.Printf("Msg Body type: %v , data: %v", ty, data)
	default:
		return fmt.Errorf("unknow type: %v", ty)
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
