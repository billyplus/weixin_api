package weixin_api

import "encoding/xml"

type BaseMessage struct {
	ToUserName   string // 开发者微信号
	FromUserName string // 发送方帐号（一个OpenID）
	CreateTime   int64  // 消息创建时间 （整型）
	MsgType      string // 消息类型
	MsgId        int64  // 消息id，64位整型
}

func DecodeRawMessage[T any](data []byte) (*T, error) {
	var raw T
	if err := xml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	return &raw, nil
}

// 文本消息 文本为text
type TextMessage struct {
	BaseMessage
	Content string // 文本消息内容
}

// ImageMessage 图片消息 图片为image
type ImageMessage struct {
	BaseMessage
	PicUrl  string // 图片链接（由系统生成）
	MediaId string // 图片消息媒体id，可以调用获取临时素材接口拉取数据。
}

// VoiceMessage 语音消息 语音为voice
type VoiceMessage struct {
	BaseMessage
	MediaId      string // 语音消息媒体id，可以调用获取临时素材接口拉取数据。
	Format       string // 语音格式，如amr，speex等
	Recongnition string // 消息id，64位整型
}

// VideoMessage 视频/小视频消息 视频为video，小视频为shortvideo
type VideoMessage struct {
	BaseMessage
	MediaId      string // 视频消息媒体id，可以调用获取临时素材接口拉取数据。
	ThumbMediaId string // 视频消息缩略图的媒体id，可以调用多媒体文件下载接口拉取数据。
}

// LocationMessage 地理位置消息，地理位置为location
type LocationMessage struct {
	BaseMessage
	LocationX float64 `xml:"Location_X"` // 地理位置纬度
	LocationY float64 `xml:"Location_Y"` // 地理位置经度
	Scale     int     // 地图缩放大小
	Label     string  // 地理位置信息
}

// LinkMessage 链接消息，链接为link
type LinkMessage struct {
	BaseMessage
	Title       string // 消息标题
	Description string // 消息描述
	Url         string // 消息链接
}
