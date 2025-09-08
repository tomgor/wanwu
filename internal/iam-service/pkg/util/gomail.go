package util

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

// MailConfig 邮件服务器配置
type MailConfig struct {
	SMTPHost  string
	SMTPPort  int
	FromEmail string
	FromName  string
	Password  string
}

// MailContent 邮件内容
type MailContent struct {
	To          []string
	Subject     string
	Body        string
	ContentType string
}

var _mailer *Mailer

// Mailer 邮件发送器
type Mailer struct {
	config MailConfig
	dialer *gomail.Dialer
}

// Init 初始化邮件发送器
func Init(smtpHost, fromEmail, password string, smtpPort int) error {
	if _mailer != nil {
		return fmt.Errorf("mailer already init")
	}
	dialer := gomail.NewDialer(smtpHost, smtpPort, fromEmail, password)
	_mailer = &Mailer{
		config: MailConfig{
			SMTPHost:  smtpHost,
			SMTPPort:  smtpPort,
			FromEmail: fromEmail,
			Password:  password,
		},
		dialer: dialer,
	}
	return nil
}

// ------------ api -----------
func SendMail(email, subject, body string) error {
	if _mailer == nil {
		return fmt.Errorf("mailer not init")
	}
	// 发送邮件
	content := MailContent{
		To:          []string{email},
		Subject:     subject,
		Body:        body,
		ContentType: "text/html",
	}
	return _mailer.send(content)
}

// ------------ internal -----------
// send 发送邮件
func (m *Mailer) send(content MailContent) error {
	// 验证必要参数
	if len(content.To) == 0 {
		return fmt.Errorf("recipient cannot be empty")
	}
	if content.Subject == "" {
		return fmt.Errorf("subject cannot be empty")
	}
	if content.Body == "" {
		return fmt.Errorf("body cannot be empty")
	}
	msg := gomail.NewMessage()
	// 设置发件人
	if m.config.FromName != "" {
		msg.SetAddressHeader("From", m.config.FromEmail, m.config.FromName)
	} else {
		msg.SetHeader("From", m.config.FromEmail)
	}
	// 设置收件人
	msg.SetHeader("To", content.To...)
	// 设置主题
	msg.SetHeader("Subject", content.Subject)
	// 设置正文类型,默认"text/html"
	contentType := "text/html"
	if content.ContentType != "" {
		contentType = content.ContentType
	}
	msg.SetBody(contentType, content.Body)
	// 发送邮件
	return m.dialer.DialAndSend(msg)
}
