package dist

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
)

var (
	MESSAGE_SUCCESS   string = "OK"
	ACCESS_KEY_ID     string = "LTAIDquqdpFdYVmW"
	ACCESS_KEY_SECRET string = "ozeRm4LWSqEnpPEw9spdOfylSNprY1"
	SIGN_NAME         string = "海拉科技"
	TEMPLATE_CODE     string = "SMS_127790203"
	PRODUCT           string = "Dysmsapi"
	VERSION           string = "2017-05-25"
	DOMAIN            string = "dysmsapi.aliyuncs.com"
	REGION            string = "cn-hangzhou"

	ErrReturnValue = errors.New("sms return message means failed.")
)

func SendSmsTo(phone string, text string) error {
	c, err := sdk.NewClientWithAccessKey(REGION, ACCESS_KEY_ID, ACCESS_KEY_SECRET)
	if nil != err {
		return err
	}
	code := fmt.Sprintf(`{"code":"%v"}`, text)

	req := requests.NewCommonRequest()
	req.Method = "POST"
	req.Scheme = "https"
	req.Domain = DOMAIN
	req.Version = VERSION
	req.ApiName = "SendSms"
	req.QueryParams["RegionId"] = REGION
	req.QueryParams["PhoneNumbers"] = phone
	req.QueryParams["SignName"] = SIGN_NAME
	req.QueryParams["TemplateCode"] = TEMPLATE_CODE
	req.QueryParams["TemplateParam"] = code

	res, err := c.ProcessCommonRequest(req)
	if nil != err {
		return err
	}
	msg := struct {
		Message   string
		RequestId string
		BizId     string
		Code      string
	}{}
	body := res.GetHttpContentBytes()
	err = json.Unmarshal(body, &msg)
	if nil != err {
		//log.Error("sms", zap.String("body", string(body)), zap.String("phone", phone))
		return err
	}
	if MESSAGE_SUCCESS != msg.Message {
		//log.Error("sms", zap.String("body", string(body)), zap.String("phone", phone))
		return ErrReturnValue
	}
	return nil
}
