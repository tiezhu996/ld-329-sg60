package errors

// BusinessError 业务异常，错误码与错误文案集中在 constants 包维护。
type BusinessError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e BusinessError) Error() string { return e.Message }

// New 按集中管理的错误码与文案构造业务异常。
func New(code, message string) BusinessError {
	return BusinessError{Code: code, Message: message}
}

// IsBusiness 判断错误是否为业务异常。
func IsBusiness(err error) bool {
	_, ok := err.(BusinessError)
	return ok
}
