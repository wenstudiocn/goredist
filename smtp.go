package dist

import (
	"github.com/go-gomail/gomail"
)

type Mailer struct {
	smtpServer string
	smtpPort int
	username string
	password string
}

type Mail struct {
	Subject string
	Body string
	Attachments []string
}

type MailReceiver struct {
	Name string
	Addr string
}

func NewMailer(smtp string, port int, username, password string) *Mailer {
	return &Mailer{
		smtpServer: smtp,
		smtpPort: port,
		username: username,
		password: password,
	}
}

func (self *Mailer) Send(mail Mail, receivers ...MailReceiver) error {
	dailer := gomail.NewDialer(self.smtpServer, self.smtpPort, self.username, self.password)
	sender, err := dailer.Dial()
	if err != nil {
		return err
	}
	msg := gomail.NewMessage()
	var ret error
	for _, receiver := range receivers {
		msg.SetHeader("From", self.username)
		msg.SetAddressHeader("To", receiver.Addr, receiver.Name)
		msg.SetHeader("Subject", mail.Subject)
		msg.SetBody("text/html", mail.Body)
		for _, file := range mail.Attachments {
			msg.Attach(file)
		}
		err = gomail.Send(sender, msg)
		if err != nil {
			ret = err
		}
		msg.Reset()
	}
	return ret
}