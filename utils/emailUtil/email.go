package emailUtil

import (
	"fmt"
	"github.com/astaxie/beego"
	"net/smtp"
	"strings"
)

/*
@Author: wangc293
@Time: 2020-05-30 10:02
@Description:
*/

const (
	Email_Name = "联通大数据智能门禁"
 	Email_Subject = "感谢您注册联通大数据智能门禁开放平台"
)

type Mail struct {
	senderId string
	toIds    []string
	subject  string
	body     string
}

type SmtpServer struct {
	host string
	port string
}

func EmailSendCode(nickname, to, code string) error {
	if nickname != "" {
		nickname = nickname + "，"
	}
	body := `
        <html>
        <body>
        <h3>
          您好：
        </h3>
        非常感谢您使用`+Email_Name+`，您的邮箱验证码为：<br/>
        <b>`+code+`</b><br/>
        此验证码有效期30分钟，请妥善保存。<br/>
        如果这不是您本人的操作，请忽略本邮件。<br/>
        </body>
        </html>
        `
	return SendToMail(to, Email_Subject, body)
}

func (s *SmtpServer) serverName() string {
	return s.host + ":" + s.port
}

func (mail *Mail) buildMessage() string {
	message := ""
	message += fmt.Sprintf("From: %s<%s>\r\n", Email_Name, mail.senderId)
	if len(mail.toIds) > 0 {
		message += fmt.Sprintf("To: %s\r\n", strings.Join(mail.toIds, ";"))
	}
	message += fmt.Sprintf("Subject: %s\r\n", mail.subject)
	message += "Content-Type: text/html; charset=UTF-8"
	message += "\r\n\r\n" + mail.body
	return message
}

func SendToMail(to, subject, body string) error {
	mail := Mail{
		senderId: getEmailUser(),
		toIds: strings.Split(to, ";"),
		subject: subject,
		body: body,
	}
	messageBody := mail.buildMessage()
	smtpServer := SmtpServer{host: getEmailHost(), port: getEmailPort()}
	//build an auth
	auth := smtp.PlainAuth("", mail.senderId, getEmailPassword(), smtpServer.host)
	sendTo := strings.Split(to, ";")
	err := smtp.SendMail(smtpServer.serverName(), auth, getEmailUser(), sendTo, []byte(messageBody))
	return err
}

func getEmailUser() string {
	return beego.AppConfig.DefaultString("email_user", "")
}

func getEmailHost() string {
	return beego.AppConfig.DefaultString("email_host", "")
}

func getEmailPort() string {
	return beego.AppConfig.DefaultString("email_port", "")
}

func getEmailPassword() string {
	return beego.AppConfig.DefaultString("email_password", "")
}