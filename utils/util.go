package utils

import (
	"html/template"
	"math/rand"
	g "odisk/global"
	"regexp"
	"time"

	"github.com/wneessen/go-mail"
)

// 检查邮箱格式
func IsValidEmail(email string) bool {
	// 邮箱正则表达式
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
/* SendEmail 发送一封邮件给指定的收件人。  
参数： 
 
- email: 收件人的邮箱地址。  

- subject: 邮件的主题。  

- t: 用于生成邮件正文的模板。  

- data: 模板所需的数据。  

返回值：  

- error: 如果在发送邮件过程中发生错误，则返回该错误；否则返回nil。
*/
func SendEmail(email string, subject string, t *template.Template, data interface{}) error {
	msg := mail.NewMsg()
	if err := msg.From(g.Config.Server.Mail.SerderMail); err != nil {
		return err
	}
	if err := msg.To(email); err != nil {
		return err
	}
	// 设置邮件的主题
	msg.Subject(subject)
	// 使用HTML模板和数据设置邮件正文
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

/*
关于函数内部的 seededRand:

seededRand 是定义的一个随机数生成器，它使用 rand.NewSource 和 rand.New 创建的。让我们深入了解一下这个机制的工作原理。
在 Go 语言的 math/rand 包中，随机数生成器是基于伪随机数生成器（PRNG）实现的。伪随机数生成器不是真正的随机，因为它们实际上是一个确定的算法，但是它们产生的序列看起来像是随机的，并且对于非加密用途来说足够好了。

rand.NewSource 函数用于创建一个新的随机数源，它接受一个 int64 类型的种子值作为参数。种子是随机数生成过程的起点。相同的种子会产生相同的随机数序列，而不同的种子通常会产生不同的序列。

rand.NewSource(time.Now().UnixNano()) 使用当前时间的纳秒数作为种子。由于时间在不断变化，因此每次调用 time.Now().UnixNano() 都会返回一个不同的值，这意味着每次创建新的随机数生成器时，你都会得到一个不同的随机数序列。

rand.New 函数接受一个 Source 接口作为参数，并返回一个 *rand.Rand 类型的新实例。这个实例可以用来生成随机数。

当你调用 seededRand.Intn(n) 时，*rand.Rand 实例会使用其内部的随机数源来生成一个介于 0 和 n-1 之间的随机整数。这个随机整数是通过一个复杂的算法从种子值派生出来的，该算法确保每次调用 Intn 时都会生成一个新的随机数，并且这些数在统计上看起来是随机的。

简而言之，seededRand 机制通过使用一个不断变化的种子值来初始化一个伪随机数生成器，然后使用该生成器来产生看似随机的数。由于种子值的变化，每次运行程序时生成的随机数序列都会不同，这对于生成验证码等一次性使用的随机值来说非常重要。
*/
func GenerateVerificationCode(length int) string {
	// 定义一个常量字符串charset，包含所有可能出现在验证码中的字符。
	const charset string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789+=-()"
	// 创建一个新的随机数生成器，使用当前时间的纳秒数作为种子，以确保每次生成的验证码都是随机的。
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	// 创建一个字节切片code，其长度等于验证码的指定长度。
	code := make([]byte, length)
	// 遍历code切片中的每个索引位置。
	for i := range code {
		// 获取字符集的长度
		charsetLength := len(charset)

		// 为当前索引位置生成一个随机字符的索引
		randomIndex := seededRand.Intn(charsetLength)

		// 从字符集中获取随机字符的字节
		randomCharByte := charset[randomIndex]

		// 将随机字符的字节存储在验证码字节切片中
		code[i] = randomCharByte
	}
	return string(code)
}
