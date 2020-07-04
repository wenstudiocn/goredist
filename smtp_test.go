package dist

import "testing"

func TestSmtp(t *testing.T) {
	mailer := NewMailer("smtp.gmail.com", 587, "", "")
	mail := Mail{Subject:"hello",Body:"good"}
	r := MailReceiver{Addr:"", Name:""}
	err := mailer.Send(mail, r)
	t.Log(err)
}