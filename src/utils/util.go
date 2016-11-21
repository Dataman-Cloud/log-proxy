package utils

import (
	"strconv"
	"time"
	"unsafe"
)

func Byte2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func ParseDate(from, to string) string {
	f, err := strconv.ParseInt(from, 10, 64)
	if err != nil {
		return time.Now().Format("2006-01-02")
	}

	t, err := strconv.ParseInt(to, 10, 64)
	if err != nil {
		return time.Now().Format("2006-01-02")
	}

	if time.Unix(f, 0).Format("2006-01-02") != time.Unix(t, 0).Format("2006-01-02") {
		return "*"
	}

	return time.Now().Format("2006-01-02")
}
