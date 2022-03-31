package weixin_api

import (
	"bytes"
	"encoding/xml"

	"github.com/pkg/errors"
)

// 微信支持的消息类型
const (
	MsgTypeText       = "text"       // 文本消息
	MsgTypeImage      = "image"      // 图片消息
	MsgTypeVoice      = "voice"      // 语音消息
	MsgTypeVideo      = "video"      // 视频消息
	MsgTypeShortVideo = "shortvideo" // 小视频消息
	MsgTypeLocation   = "location"   // 地理位置消息
	MsgTypeLink       = "link"       // 链接消息
	MsgTypeEvent      = "event"      // 事件推送
)

var (
	ErrInvalidHandler     = errors.New("无效的消息处理函数")
	ErrInvalidMessageType = errors.New("无效的消息类型")
)

func (e *Engine) HandleMessage(body []byte) error {
	decoder := xml.NewDecoder(bytes.NewBuffer(body))
	msgTyp := ""
	var err error

	for {
		// Read tokens from the XML document in a stream.
		t, _ := decoder.Token()
		if t == nil {
			err = errors.New("xml解析token出错")
			break
		}
		// Inspect the type of the token just read.
		switch se := t.(type) {
		case xml.StartElement:
			// 找到消息类型
			if se.Name.Local == "MsgType" {
				// 解析消息类型
				err = decoder.DecodeElement(&msgTyp, &se)
				break
			}
		}
	}

	if err != nil {
		return errors.Wrap(err, "DecodeXML")
	}

	switch msgTyp {
	case MsgTypeEvent:
		return handle(e.handleTextMessage, body)
	case MsgTypeImage:
		return handle(e.handleImageMessage, body)
	case MsgTypeVoice:
		return handle(e.handleVoiceMessage, body)
	case MsgTypeVideo:
		return handle(e.handleVideoMessage, body)
	case MsgTypeLocation:
		return handle(e.handleLocationMessage, body)
	case MsgTypeLink:
		return handle(e.handleLinkMessage, body)
	default:
		return ErrInvalidMessageType
	}
}

func handle[T any](fn func(m *T) error, body []byte) error {
	if fn == nil {
		return ErrInvalidHandler
	}
	m, err := DecodeRawMessage[T](body)
	if err != nil {
		return err
	}
	return fn(m)
}
