package utils

import (
	"unsafe"
)

func Byte2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
