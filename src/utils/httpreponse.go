package utils

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	CODE_OK       = 0
	CODE_UNDEFINE = 10001
)

type Error struct {
	Code string `json:"code"`
	Err  error  `json:"Err"`
}

func NewError(code string, err error) *Error {
	return &Error{
		Code: code,
		Err:  err,
	}
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func Ok(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{"code": CODE_OK, "data": data})
}

func Create(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusCreated, gin.H{"code": CODE_OK, "data": data})
	return
}

func Delete(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusNoContent, gin.H{"code": CODE_OK, "data": data})
	return
}

func Update(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusAccepted, gin.H{"code": CODE_OK, "data": data})
	return
}

func ErrorResponse(ctx *gin.Context, err error) {
	hcode := http.StatusServiceUnavailable
	ecode := CODE_UNDEFINE

	e, ok := err.(*Error)
	fmt.Println(e, ok)
	if ok {
		if codes := strings.Split(e.Code, "-"); len(codes) == 2 {
			hcode, _ = strconv.Atoi(codes[0])
			ecode, _ = strconv.Atoi(codes[1])
			ctx.JSON(hcode, gin.H{"code": ecode, "data": err.Error()})
			return
		}
	}
	ctx.JSON(hcode, gin.H{"code": ecode, "data": err.Error()})
}
