package utils

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	// CodeOK ok status
	CodeOK = 0
	// CodeUndefine undefine status
	CodeUndefine = 10001
)

// Error struct
type Error struct {
	Code string `json:"code"`
	Err  error  `json:"Err"`
}

// NewError new error
func NewError(code string, err interface{}) *Error {
	e, ok := err.(error)
	if !ok {
		e = errors.New(fmt.Sprint(err))
	}
	return &Error{
		Code: code,
		Err:  e,
	}
}

// Error parse error to string
func (e *Error) Error() string {
	return e.Err.Error()
}

// Ok return 200 status
func Ok(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{"code": CodeOK, "data": data})
}

// Create return create status
func Create(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusCreated, gin.H{"code": CodeOK, "data": data})
	return
}

// Delete return delete status
func Delete(ctx *gin.Context, data interface{}) {
	//ctx.JSON(http.StatusNoContent, gin.H{"code": CodeOK, "data": data})
	ctx.Writer.WriteHeader(http.StatusNoContent)
	return
}

// Update return update status
func Update(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusAccepted, gin.H{"code": CodeOK, "data": data})
	return
}

// ErrorResponse return error status
func ErrorResponse(ctx *gin.Context, err error) {
	hcode := http.StatusServiceUnavailable
	ecode := CodeUndefine

	e, ok := err.(*Error)
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
