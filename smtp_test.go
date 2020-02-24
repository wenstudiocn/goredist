package dist

import "testing"

func TestSmtp(t *testing.T) {
	mailer := NewMailer("smtp.gmail.com", 587, "", "")
	r := MailReceiver{Mail:Mail{Subject:"hello",Body:"good"}, Addr:"", Name:""}
	err := mailer.Send(r)
	t.Log(err)
}