package utils

import (
	"time"
	"unsafe"
)

func Byte2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func ParseDate(from, to interface{}) string {
	f, ok := from.(int64)
	if !ok {
		return time.Now().Format("2006-01-02")
	}

	t, ok := to.(int64)
	if !ok {
		return time.Now().Format("2006-01-02")
	}

	if time.Unix(f/1000, 0).Format("2006-01-02") != time.Unix(t/1000, 0).Format("2006-01-02") {
		return "*"
	}

	return time.Now().Format("2006-01-02")
}
