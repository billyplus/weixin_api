package weixin_api

import "fmt"

type ErrorMsg struct {
	ErrCode int32  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func (err *ErrorMsg) Error() string {
	return fmt.Sprintf("%d: %s", err.ErrCode, err.ErrMsg)
}
