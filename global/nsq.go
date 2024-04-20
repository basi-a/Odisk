package global

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"github.com/nsqio/go-nsq"
	"gopkg.in/gomail.v2"
)
type EmailData struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}
type Message struct {
	MsgType    string `json:"msgtype"` // 消息类型（例如，“email”、“simple”）
	DataBase64 string `json:"data"`
}

type SendEmailData struct {
	Email      string `json:"email"`
	Subject    string `json:"subject"`
	DataBase64 string `json:"data"`
}

type SimpleData struct {
	Str string `json:"str"`
}

var Producer *nsq.Producer
var Consumers map[int]*nsq.Consumer

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
	Producer.SetLoggerLevel(nsq.LogLevelWarning)
	if err = Producer.Ping(); err != nil {
		log.Println("Error pinging nsqd:", err)
		Producer.Stop()
		return err
	}
	simple := SimpleData{
		Str: "topic created",
	}
	jsonData, _ := json.Marshal(simple)
	for topic := range Config.Nsq.Topics {
		if err = ProduceMsg(topic, "simple", jsonData); err != nil {
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

	// consumers := make([]*nsq.Consumer, 0)
	consumers := make(map[int]*nsq.Consumer, 0)
	count := -1
	for topic, channals := range Config.Nsq.Topics {

		for _, channal := range channals {
			consumer, err := nsq.NewConsumer(topic, channal, nsq.NewConfig())
			if err != nil {
				return err
			}
			consumer.SetLoggerLevel(nsq.LogLevelWarning)

			consumer.AddConcurrentHandlers(nsq.HandlerFunc(ConsumeMsg), 5)

			if err = consumer.ConnectToNSQLookupds(nsqlookupdAddrsWithPort); err != nil {
				return err
			}
			count++
			// consumers = append(consumers, consumer)
			consumers[count]=consumer
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
		data, err := base64.RawStdEncoding.DecodeString(msg.DataBase64)
		if err != nil {
			return nil
		}
		var sendEmailData SendEmailData
		if err := json.Unmarshal(data, &sendEmailData); err != nil {
			return fmt.Errorf("unmarshal err: %s \n| %s", err, string(data))
		}

		if err := sendEmail(sendEmailData.Email, sendEmailData.Subject, sendEmailData.DataBase64); err != nil {
			return fmt.Errorf("send mail error:  %s\n| %s", err, string(data))
		}

		m.Finish()
	case "simple":
		data, err := base64.RawStdEncoding.DecodeString(msg.DataBase64)
		if err != nil {
			return nil
		}
		var simpleData SimpleData
		if err := json.Unmarshal(data, &simpleData); err != nil {
			return err
		}
		fmt.Println(simpleData.Str)

		m.Finish()
	default:
		return fmt.Errorf("unknown msg type: %s", msg.MsgType)
	}
	return nil
}

func ProduceMsg(topic string, msgType string, data []byte) error {

	msg := Message{
		MsgType:    msgType,
		DataBase64: base64.RawStdEncoding.EncodeToString(data),
	}
	jsonData, err := json.Marshal(msg)

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

func sendEmail(email string, subject string, base64data string) error {
	config := Config.Server.Mail
	jsonData, err := base64.RawStdEncoding.DecodeString(base64data)
	if err != nil {
		return err
	}
	// return fmt.Errorf(string(jsonData))
	var htmlContent bytes.Buffer

	data := new(EmailData)
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return err
	}
	if err := EmailTemplate.Execute(&htmlContent, data); err != nil {
		return err
	}
	msg := gomail.NewMessage()
	msg.SetHeader("From", config.SenderMail)
	msg.SetHeader("To", email)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", htmlContent.String())

	dialer := gomail.NewDialer(config.SmtpServer, config.Port, config.UserName, config.Password)
	if err := dialer.DialAndSend(msg); err != nil {
		return fmt.Errorf("DialAndSend err: %s", err)
	}

	return nil
}
