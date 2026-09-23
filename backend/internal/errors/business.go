package errors

// BusinessError 业务异常：错误码、提示信息和对应的 HTTP 状态码集中维护
type BusinessError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
}

func (e BusinessError) Error() string { return e.Message }

// New 构造业务异常
func New(status int, code, message string) BusinessError {
	return BusinessError{Code: code, Message: message, Status: status}
}
