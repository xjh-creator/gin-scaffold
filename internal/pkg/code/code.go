package code

import (
	"fmt"
	"net/http"
)

type Error struct {
	// 错误码
	code int `json:"code"`
	// 错误消息
	msg string `json:"msg"`
	// 详细信息
	details []string `json:"details"`
}

var Codes = map[int]string{
	SUCCESS: "成功",
	ERROR:   "失败",

	ServerError:               "服务内部错误",
	InvalidParams:             "入参错误",
	NotFound:                  "找不到",
	UnauthorizedAuthNotExist:  "鉴权失败，找不到对应的AppKey和AppSecret",
	UnauthorizedTokenError:    "鉴权失败，Token错误",
	UnauthorizedTokenTimeout:  "鉴权失败，Token超时",
	UnauthorizedTokenGenerate: "鉴权失败，Token生成失败",
	TooManyRequests:           "请求过多",
}

//GetMsg get error information based on Code
func GetMsg(code int) string {
	msg, ok := Codes[code]
	if ok {
		return msg
	}
	return Codes[ERROR]
}

func NewError(code int, msg string) *Error {
	if _, ok := Codes[code]; ok {
		panic(fmt.Sprintf("错误码 %d 已经存在，请更换一个", code))
	}
	Codes[code] = msg
	return &Error{code: code, msg: msg}
}

func (e *Error) Error() string {
	return fmt.Sprintf("错误码：%d, 错误信息:：%s", e.Code(), e.Msg())
}

func (e *Error) Code() int {
	return e.code
}

func (e *Error) Msg() string {
	return e.msg
}

func (e *Error) Msgf(args []interface{}) string {
	return fmt.Sprintf(e.msg, args...)
}

func (e *Error) Details() []string {
	return e.details
}

func (e *Error) WithDetails(details ...string) *Error {
	newError := *e
	newError.details = []string{}
	for _, d := range details {
		newError.details = append(newError.details, d)
	}

	return &newError
}

func (e *Error) StatusCode() int {
	switch e.Code() {
	case SUCCESS:
		return http.StatusOK
	case ServerError:
		return http.StatusInternalServerError
	case InvalidParams:
		return http.StatusBadRequest
	case UnauthorizedAuthNotExist:
		fallthrough
	case UnauthorizedTokenError:
		fallthrough
	case UnauthorizedTokenGenerate:
		fallthrough
	case UnauthorizedTokenTimeout:
		return http.StatusUnauthorized
	case TooManyRequests:
		return http.StatusTooManyRequests
	}

	return http.StatusInternalServerError
}
