package global

import (

	"encoding/json"
	"fmt"
	"log"

	"html/template"

	"github.com/nsqio/go-nsq"
	"github.com/wneessen/go-mail"
)

type Message struct {
	MsgType string      `json:"msgtype"` // 消息类型（例如，“email”、“simple”）
	Data    interface{} `json:"data"`    // 消息数据（根据类型而变化）
}

type EmailData struct {
	Email   string             `json:"email"`
	Subject string             `json:"subject"`
	T       *template.Template `json:"t"`
	Data    interface{}        `json:"data"`
}

type SimpleData struct {
	Str string `json:"str"`
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
	simple := SimpleData{
		Str: "topic created",
	}

	for topic := range Config.Nsq.Topics {
		if err = ProduceMsg(topic, "simple", simple); err != nil {
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

			consumer.AddConcurrentHandlers(nsq.HandlerFunc(ConsumeMsg), 2)

			if err = consumer.ConnectToNSQLookupds(nsqlookupdAddrsWithPort); err != nil {
				return err
			}

			consumers = append(consumers, consumer)
		}
	}
	Consumers = consumers
	return nil
}

func ConsumeMsg(m *nsq.Message) error {
	// log.Println(string(m.Body))
	msg := new(Message)
	if err := json.Unmarshal(m.Body, &msg); err != nil {
		return err
	}
	switch msg.MsgType {
	case "email":
		 if emailData, ok := msg.Data.(EmailData); ok {
			if err := sendEmail(emailData.Email, emailData.Subject, emailData.T, emailData.Data); err != nil {
				return fmt.Errorf("send mail error: %v", err)
			}
		 }else {
			return fmt.Errorf("type assertion failed")
		}
	case "simple":
		if simpledata, ok := msg.Data.(SimpleData); ok {
			log.Printf("simple msg: %s", simpledata.Str)
		}else {
			return fmt.Errorf("type assertion failed")
		}
	default:
		return fmt.Errorf("unknown msg type: %s", msg.MsgType)
	}
	return nil
}

func ProduceMsg(topic string, msgType string, data interface{}) error {
	msg := Message{
		MsgType: msgType,
		Data:    data,
	}
	jsonData, err := json.Marshal(msg)
	// log.Println(string(jsonData))
	if err != nil {
		return err
	}
	if err := Producer.Publish(topic, jsonData); err != nil {
		return err
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
