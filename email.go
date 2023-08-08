package utils

import (
	"errors"
	"fmt"
	"gopkg.in/gomail.v2"
	"mime"
	"os"
)

type Email struct {
	option   server          // 基础配置
	Attach   MailAttach      // 邮件附件
	NickName string          // 发送方别名
	HookFunc EmailHandleFunc // 自定义回调函数

}

// EmailHandleFunc 自定义回调函数，如果发送的消息无法满足你的需求你可以自己配置。
//
// 例如：设置邮件正文 	func (m *gomail.Message){ m.SetBody("text/html", body) }
type EmailHandleFunc func(m *gomail.Message)

// server 基本配置信息
type server struct {
	user    string
	pass    string
	host    string
	port    int
	message *gomail.Message
}

type MailAttach struct {
	carry    bool
	nickName string
	filePath string
}

// NewEmail 构建邮件结构实体对象。
//
//	user:发送人邮箱（邮箱以自己的为准） pass: 发送人邮箱的授权码，现在可能会需要邮箱 开启授权密码后在pass填写授权码
//	host:邮箱服务器（此时用的是qq邮箱） port: 邮件服务所需端口，默认465.
func NewEmail(user, pass, host string) *Email {
	email := Email{
		option: server{
			user: user,
			pass: pass,
			host: host,
			port: 465,
		},
		NickName: "go-utils-Email",
	}
	email.option.message = email.defaultMes()
	return &email
}

func (e *Email) defaultMes() *gomail.Message {
	msg := gomail.NewMessage(
		gomail.SetCharset("UTF-8"),
	)
	return msg
}

// SetMsgOpt 如果默认的邮件格式不符合你的预期，你可以重新制定。
func (e *Email) SetMsgOpt(opts ...gomail.MessageSetting) {
	e.option.message = gomail.NewMessage(opts...)
}

// SendAttach  发送附件
func (e *Email) SendAttach(filePath string, nickName string) error {
	if len(filePath) == 0 {
		return errors.New("文件路径不存在")
	}
	e.Attach.carry = true
	e.Attach.filePath = filePath
	e.Attach.nickName = nickName
	return nil
}

func (e *Email) SetPort(port int) {
	e.option.port = port
}

// SendMail 发送邮件
func (e *Email) SendMail(mailTo []string, subject string, body string) error {
	// 设置邮箱主题：默认格式
	e.option.message.SetHeader("From", e.option.message.FormatAddress(e.option.user, e.NickName)) // 添加别名
	e.option.message.SetHeader("To", mailTo...)                                                   // 发送给用户(可以多个)
	e.option.message.SetHeader("Subject", subject)                                                // 设置邮件主题
	e.option.message.SetBody("text/html", body)                                                   // 设置邮件正文
	// 触发回调，加载用户自定义的格式设置
	e.HookFunc(e.option.message)

	if e.Attach.carry {
		e.option.message.Attach(e.Attach.filePath,
			gomail.Rename(e.Attach.nickName), //重命名
			gomail.SetHeader(map[string][]string{
				"Content-Disposition": {
					fmt.Sprintf(`attachment; filename="%s"`, mime.QEncoding.Encode("UTF-8", e.Attach.nickName)),
				},
			}),
		)
	}

	/*
	   创建SMTP客户端，连接到远程的邮件服务器，需要指定服务器地址、端口号、用户名、密码，如果端口号为465的话，
	   自动开启SSL，这个时候需要指定TLSConfig
	*/

	d := gomail.NewDialer(e.option.host, e.option.port, e.option.user, e.option.pass) // 设置邮件正文
	// TODO 暂不支持使用TLS
	//d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	err := d.DialAndSend(e.option.message)
	return err
}

// ParseString 可以直接读取HTML并保存为Byte,也就是发送邮件可以直接发送一个内嵌的HTML.
func (e *Email) ParseString(path string) string {
	byte, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("ReadFile error:", err)
		return ""
	}
	return string(byte)
}
