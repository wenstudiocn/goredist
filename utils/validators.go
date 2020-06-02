package utils

import "regexp"

func ValidPhoneCN(phone string) bool {
	//regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	pattern := "^1[3-9]\\d{9}$"
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(phone)
}