package utils

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

var (
	s         = rand.NewSource(time.Now().UnixNano())
	r         = rand.New(s)
	CHAR_POOL = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
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
