package weixin_api

import (
	"bytes"
	"context"
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

const (
	EventTypeClick       = "CLICK"       // 自定义菜单事件
	EventTypeView        = "VIEW"        // 点菜单跳转链接
	EventTypeLocation    = "LOCATION"    // 上报地理位置
	EventTypeScan        = "SCAN"        // 用户已关注
	EventTypeSubscribe   = "subscribe"   // 用户未关注
	EventTypeUnsubscribe = "unsubscribe" // 取消订阅

)

var (
	ErrInvalidHandler = errors.New("无效的消息处理函数")
)

func (e *Engine) HandleMessage(c context.Context, data []byte) error {
	decoder := xml.NewDecoder(bytes.NewBuffer(data))
	msgTyp := ""
	evTyp := ""
	var err error

	for {
		// TODO：快速查找指定节点
		t, err := decoder.Token()
		if err != nil {
			break
		}
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
	case MsgTypeText:
		return handle(e.handleTextMessage, data)
	case MsgTypeImage:
		return handle(e.handleImageMessage, data)
	case MsgTypeVoice:
		return handle(e.handleVoiceMessage, data)
	case MsgTypeVideo:
		return handle(e.handleVideoMessage, data)
	case MsgTypeLocation:
		return handle(e.handleLocationMessage, data)
	case MsgTypeLink:
		return handle(e.handleLinkMessage, data)
	case MsgTypeEvent:
		for {
			// TODO：快速查找指定节点
			t, err := decoder.Token()
			if err != nil {
				break
			}
			if t == nil {
				err = errors.New("xml解析token出错")
				break
			}
			// Inspect the type of the token just read.
			switch se := t.(type) {
			case xml.StartElement:
				// 找到消息类型
				if se.Name.Local == "Event" {
					// 解析消息类型
					err = decoder.DecodeElement(&evTyp, &se)
					break
				}
			}
		}

		if err != nil {
			return errors.Wrap(err, "DecodeXML")
		}

		switch evTyp {
		case EventTypeSubscribe:
			return handle(e.handleSubscribeEvent, data)
		case EventTypeUnsubscribe:
			return handle(e.handleUnsubscribeEvent, data)
		case EventTypeClick:
			return handle(e.handleClickEvent, data)
		case EventTypeView:
			return handle(e.handleViewEvent, data)
		case EventTypeLocation:
			return handle(e.handleLocationEvent, data)
		case EventTypeScan:
			return handle(e.handleScanEvent, data)
		}
		return &ErrInvalidEventType{Type: evTyp}
	}

	return &ErrInvalidMessageType{Type: msgTyp}
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
