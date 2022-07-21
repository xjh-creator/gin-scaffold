package code

const (
	SUCCESS = 200 //成功
	ERROR   = 404 //失败

	ServerError = iota + 10000000
	InvalidParams
	NotFound
	UnauthorizedAuthNotExist
	UnauthorizedTokenError
	UnauthorizedTokenTimeout
	UnauthorizedTokenGenerate
	TooManyRequests
)
