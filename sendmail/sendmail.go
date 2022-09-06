package sendmail

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"strings"
)

// SendMail 发邮件
func SendMail(mailTo string, subject string, body string) error {
	m := gomail.NewMessage()
	//设置发件人
	m.SetHeader("From", "444216978@qq.com")
	//设置发送给多个用户
	mailArrTo := strings.Split(mailTo, ",")
	m.SetHeader("To", mailArrTo...)
	//设置邮件主题
	m.SetHeader("Subject", subject)
	//设置邮件正文
	m.SetBody("text/html", body)
	d := gomail.NewDialer("email.qq.com", 587, "444216978@qq.com", "oxlqpamltwqibieg")
	err := d.DialAndSend(m)
	if err != nil {
		fmt.Println(err)
	}
	return err
}
