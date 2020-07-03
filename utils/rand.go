package utils

import (
	"fmt"
	"github.com/mojocn/base64Captcha"
	"math"
	"math/rand"
	"time"
)

var (
	s         = rand.NewSource(time.Now().UnixNano())
	r         = rand.New(s)
	CHAR_POOL = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	captcha_digit_driver = base64Captcha.NewDriverDigit(80, 200, 5, 0.7, 80)
	captcha_store        = base64Captcha.NewMemoryStore(64, time.Second*100)
	captchaDigit         = base64Captcha.NewCaptcha(captcha_digit_driver, captcha_store)
)

func RandInt(max int) int {
	return r.Intn(max)
}

func RandIntScope(from, to int) int {
	return r.Intn(to-from) + from
}

func RandNumString(length int) string {
	max := int32(math.Pow10(length))
	n := r.Int31n(max)
	sf := fmt.Sprintf("%%0%dv", length)
	return fmt.Sprintf(sf, n)
}

func RandNumAlphaString(length int) string {
	l := len(CHAR_POOL)
	str := ""

	for i := 0; i < length; i++ {
		str += fmt.Sprintf("%c", CHAR_POOL[r.Intn(l)])
	}
	return str
}

func DigitCaptchaGen() (string, string, error) {
	return captchaDigit.Generate()
}

func DigitCaptchaVerify(id string, answer string) bool {
	return captchaDigit.Verify(id, answer, false)
}
