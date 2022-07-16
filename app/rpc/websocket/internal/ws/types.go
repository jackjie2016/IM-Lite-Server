package ws

import "fmt"

var (
	ErrMaxConnNum = fmt.Errorf("max conn num")
)

type ConnRequest struct {
	UserID   string `form:"userID"`
	Platform string `form:"platform"`
	Token    string `form:"token"`
}

type ConnResponse struct {
	Uid     string `json:"uid"`
	ErrMsg  string `json:"errMsg"`
	Success bool   `json:"success"`
}
