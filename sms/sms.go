package sms

import "errors"

const (
	http_status_code_err_tpl = "bad status code %d"
)

var (
	ErrSmsBadPlatformConfig = errors.New("error sms platform config")
	ErrSmsNoMoney = errors.New("sms no money")
	ErrSmsReturnBadValue = errors.New("sms return error value")
)

type ISms interface {
	SendVCode(phone, code string) error
	GetBalance() (float64, error)
	GetRemains() (int, error)
}

func NewSmsBao(username, password, tplSmsVcode string) ISms {
	return &smsBao{
		username: username,
		password: password,
		tplSmsVcode: tplSmsVcode,
	}
}

func NewSmsFastoo() ISms {
	return &smsFastoo{}
}
