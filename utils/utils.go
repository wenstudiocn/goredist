package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/wenstudiocn/goredist/e"
	"os"
	"strings"
	"unsafe"
)

const (
	TIME_FORMAT = "2006-01-02 15:04:05"
)

var (
	ErrPathIsNotAPath = errors.New("Not a path")
	ErrPathPerm       = errors.New("User has not perm to path")
)

func ErrPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func Unused(v ...interface{}) {

}

func Md5(data []byte) string {
	dest := md5.Sum(data)
	return fmt.Sprintf("%x", dest)
}

func Sha1(data []byte) string {
	dest := sha1.Sum(data)
	return fmt.Sprintf("%x", dest)
}

func Sha256(data []byte) string {
	dest := sha256.Sum256(data)
	return fmt.Sprintf("%x", dest)
}

func MemSet(p unsafe.Pointer, b byte, length uintptr) {
	np := uintptr(p)
	var i uintptr
	for i = 0; i < length; i++ {
		pb := (*byte)(unsafe.Pointer(np + i))
		*pb = b
	}
}

// 配置中的路径存在检查, 自动创建, 并执行权限检查
func EnsurePath(path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
	} else {
		if !fi.IsDir() {
			return ErrPathIsNotAPath
		}
		mod := fi.Mode()
		if mod&0700 == 0 {
			return ErrPathPerm
		}
	}
	return nil
}
// parse name from an email address
func GetNameFromEmail(email string) (string, error) {
	if ValidEmail(email) {
		parts := strings.Split(email, "@")
		return parts[0], nil
	}
	return email, e.ErrParams
}
