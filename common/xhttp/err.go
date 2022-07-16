package xhttp

type XResponse struct {
	Code int32       `json:"C"`
	Msg  string      `json:"M"`
	Data interface{} `json:"data,omitempty"`
}

type ICodeErr interface {
	Code() int32
	Error() string
}

func Success(data interface{}) *XResponse {
	return &XResponse{Code: 0, Data: data}
}

func Failed(err ICodeErr) *XResponse {
	return &XResponse{Code: err.Code(), Msg: err.Error()}
}

type CodeError struct {
	C int32  `json:"code"`
	M string `json:"msg"`
}

func (c *CodeError) Code() int32 {
	return c.C
}

func (c *CodeError) Error() string {
	return c.M
}

func NewErr(code int32, msg string) *CodeError {
	return &CodeError{C: code, M: msg}
}

func NewDefaultErr(msg string) *CodeError {
	return NewErr(defaultErrCode, msg)
}

func NewParamErr(err error) *CodeError {
	return NewErr(paramErrCode, err.Error())
}

func NewParamErrByMsg(msg string) *CodeError {
	return NewErr(paramErrCode, msg)
}

var internalErr = &CodeError{C: internalErrCode, M: "服务繁忙，请稍后再试"}

func NewInternalErr() *CodeError {
	return internalErr
}
