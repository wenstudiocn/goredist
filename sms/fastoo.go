package sms

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

/* fastoo.cn */

const (
	SMS_PLATFORM_FASTOO = "smsFastoo"
	fastoo_account_url = "http://api.fastoo.cn/v1/admin/getUserAccounts.json"
	fastoo_msg_url     = "http://api.fastoo.cn/v1/submit.json"
	fastoo_api_key     = "a47dd09338834deda2e0e717a90e2cce"

	fastoo_err_code_tpl = "fastoo return errcode %d"
)

type smsFastoo struct{}

type fastooSmsJsonParam struct {
	ApiKey string `json:"apiKey"`
	Da     string `json:"da"`
	Msg    string `json:"msg"`
}

type fastooSmsJsonResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func newFastooSmsJsonParam(phone, msg string) *fastooSmsJsonParam {
	return &fastooSmsJsonParam{
		ApiKey: fastoo_api_key,
		Da:     "86" + phone,
		Msg:    msg,
	}
}

func (self *smsFastoo) SendVCode(phone, code string) error {
	//msg := "hello" //fmt.Sprintf(tpl, "YQH", code)
	form := url.Values{}
	form.Add("apiKey", fastoo_api_key)
	//form.Add("da", "8618638561215")
	//form.Add("msg", msg)
	//dest := "apiKey=" + fastoo_api_key + "&da=86" + phone + "&msg=" + msg
	fmt.Println(form.Encode())
	resp, err := http.PostForm(fastoo_account_url, form)
	//resp, err := http.Post(fastoo_url, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		r := fastooSmsJsonResult{}
		err = json.Unmarshal(body, &r)
		if err != nil {
			return err
		}
		if r.Code != 0 {
			return errors.New(fmt.Sprintf(fastoo_err_code_tpl, r.Code))
		}
	} else {
		return errors.New(fmt.Sprintf(http_status_code_err_tpl, resp.StatusCode))
	}

	return nil
}

func (self *smsFastoo) GetBalance() (float64, error) {
	return 0.0, nil
}

func (self *smsFastoo) GetRemains() (int, error) {
	return 0, nil
}