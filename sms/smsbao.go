package sms

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	SMS_PLATFORM_SMSBAO = "smsBao"
	url_SENDSMS = "http://api.smsbao.com/sms"
	url_BALANCE = "http://www.smsbao.com/query?u=%s&p=%s"
	query = "u=%s&p=%s&m=%s&c=%s"
)

type smsBao struct {
	tplSmsVcode string
	username, password string
}

func (self *smsBao) SendVCode(phone, vcode string) error {
	content := fmt.Sprintf(self.tplSmsVcode, vcode)
	uri := fmt.Sprintf("%v?"+query, url_SENDSMS, self.username, self.password, phone, url.QueryEscape(content))
	resp, err := http.Get(uri)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		sbody := fmt.Sprintf("%s", body)
		if sbody != "0" {
			return ErrSmsNoMoney
		}
	} else {
		return errors.New(fmt.Sprintf(http_status_code_err_tpl, resp.StatusCode))
	}

	return nil
}

func (self *smsBao) GetBalance() (float64, error) {
	return 0.0, nil
}

func (self *smsBao) GetRemains() (int, error) {
	uri := fmt.Sprintf(url_BALANCE, self.username, self.password)
	resp, err := http.Get(uri)
	if err != nil {
		return 0.0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return 0, err
		}
		sbody := fmt.Sprintf("%s", body)

		lines := strings.Split(sbody, "\n")
		if len(lines) != 2 {
			return 0, ErrSmsReturnBadValue
		}
		if lines[0] != "0" {
			return 0, ErrSmsReturnBadValue
		}
		parts := strings.Split(lines[1], ",")
		if len(parts) != 2 {
			return 0, ErrSmsReturnBadValue
		}
		remains, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return 0.0, err
		}
		return int(remains), nil
	} else {
		return 0, errors.New(fmt.Sprintf(http_status_code_err_tpl, resp.StatusCode))
	}
}