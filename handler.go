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
	ErrInvalidHandler  = errors.New("未注册消息处理函数")
	ErrInvalidXMLToken = errors.New("xml解析token出错")
)

type TextMessageHandler func(m *TextMessage) error
type ImageMessageHandler func(m *ImageMessage) error
type VoiceMessageHandler func(m *VoiceMessage) error
type VideoMessageHandler func(m *VideoMessage) error
type LocationMessageHandler func(m *LocationMessage) error
type LinkMessageHandler func(m *LinkMessage) error
type ClickEventHandler func(m *ClickEvent) error
type ViewEventHandler func(m *ViewEvent) error
type LocationEventHandler func(m *LocationEvent) error
type ScanEventHandler func(m *ScanEvent) error
type SubscribeEventHandler func(m *SubscribeEvent) error
type UnsubscribeEventHandler func(m *UnsubscribeEvent) error

func (e *Engine) HandleMessage(c context.Context, data []byte) error {
	decoder := xml.NewDecoder(bytes.NewBuffer(data))
	msgTyp := ""
	evTyp := ""
	var err error
	var t xml.Token

LOOP:
	for {
		// TODO：快速查找指定节点
		t, err = decoder.Token()
		if err != nil {
			break LOOP
		}
		if t == nil {
			err = ErrInvalidXMLToken
			break LOOP
		}
		// Inspect the type of the token just read.
		switch se := t.(type) {
		case xml.StartElement:
			// 找到消息类型
			if se.Name.Local == "MsgType" {
				// 解析消息类型
				t, err = decoder.Token()
				if err != nil {
					break LOOP
				}
				if t == nil {
					err = ErrInvalidXMLToken
					break LOOP
				}
				switch el := t.(type) {
				case xml.CharData:
					msgTyp = string(el)
				}
				// err = decoder.DecodeElement(&msgTyp, &se)
				break LOOP
			}
			break
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
	LOOP_EVENT:
		for {
			// TODO：快速查找指定节点
			t, err := decoder.Token()
			if err != nil {
				break
			}
			if t == nil {
				err = ErrInvalidXMLToken
				break
			}
			// Inspect the type of the token just read.
			switch se := t.(type) {
			case xml.StartElement:
				// 找到消息类型
				if se.Name.Local == "Event" {
					// 解析event类型
					t, err = decoder.Token()
					if err != nil {
						break LOOP_EVENT
					}
					if t == nil {
						err = ErrInvalidXMLToken
						break LOOP_EVENT
					}
					switch el := t.(type) {
					case xml.CharData:
						evTyp = string(el)
					}
					// err = decoder.DecodeElement(&msgTyp, &se)
					break LOOP_EVENT
				}
				break
			}
		}

		if err != nil {
			return errors.WithMessage(err, "DecodeMessageType")
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

func (e *Engine) RegTextMessageHandler(h TextMessageHandler) {
	e.handleTextMessage = h
}

func (e *Engine) RegImageMessageHandler(h ImageMessageHandler) {
	e.handleImageMessage = h
}

func (e *Engine) RegVoiceMessageHandler(h VoiceMessageHandler) {
	e.handleVoiceMessage = h
}

func (e *Engine) RegVideoMessageHandler(h VideoMessageHandler) {
	e.handleVideoMessage = h
}

func (e *Engine) RegLocationMessageHandler(h LocationMessageHandler) {
	e.handleLocationMessage = h
}

func (e *Engine) RegLinkMessageHandler(h LinkMessageHandler) {
	e.handleLinkMessage = h
}

func (e *Engine) RegClickEventHandler(h ClickEventHandler) {
	e.handleClickEvent = h
}

func (e *Engine) RegViewEventHandler(h ViewEventHandler) {
	e.handleViewEvent = h
}

func (e *Engine) RegLocationEventHandler(h LocationEventHandler) {
	e.handleLocationEvent = h
}

func (e *Engine) RegScanEventHandler(h ScanEventHandler) {
	e.handleScanEvent = h
}

func (e *Engine) RegSubscribeEventHandler(h SubscribeEventHandler) {
	e.handleSubscribeEvent = h
}

func (e *Engine) RegUnsubscribeEventHandler(h UnsubscribeEventHandler) {
	e.handleUnsubscribeEvent = h
}
