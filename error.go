package weixin_api

import "fmt"

type ErrorMsg struct {
	ErrCode int32  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func (err *ErrorMsg) Error() string {
	return fmt.Sprintf("%d: %s", err.ErrCode, err.ErrMsg)
}

type ErrInvalidMessageType struct {
	Type string
}

func (err *ErrInvalidMessageType) Error() string {
	return fmt.Sprintf("无效的消息类型:%s", err.Type)
}

type ErrInvalidEventType struct {
	Type string
}

func (err *ErrInvalidEventType) Error() string {
	return fmt.Sprintf("无效的事件类型:%s", err.Type)
}
