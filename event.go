package weixin_api

type BaseEvent struct {
	BaseMessage
	Event string // 事件类型，subscribe(订阅)、unsubscribe(取消订阅)
}

// 上报地理位置事件 事件类型，LOCATION
type LocationEvent struct {
	BaseEvent
	Latitude  float32 //地理位置纬度
	Longitude float32 //地理位置经度
	Precision float32 //地理位置精度
}

// 自定义菜单事件 事件类型，CLICK
type ClickEvent struct {
	BaseEvent
	EventKey string // 事件KEY值，与自定义菜单接口中KEY值对应
}

// 点击菜单跳转链接时的事件推送 事件类型，VIEW
type ViewEvent struct {
	BaseEvent
	EventKey string // 事件KEY值，设置的跳转URL

}

// 用户已关注时的事件推送
type ScanEvent struct {
	BaseEvent
	EventKey string // 事件KEY值，是一个32位无符号整数，即创建二维码时的二维码scene_id
	Ticket   string //二维码的ticket，可用来换取二维码图片
}

// 事件类型，subscribe 事件类型，unsubscribe(取消订阅)
type SubscribeEvent struct {
	BaseEvent
	EventKey string // 事件KEY值，qrscene_为前缀，后面为二维码的参数值
	Ticket   string // 二维码的ticket，可用来换取二维码图片
}

// 事件类型，unsubscribe 事件类型，unsubscribe(取消订阅)
type UnsubscribeEvent struct {
	BaseEvent
	EventKey string // 事件KEY值，qrscene_为前缀，后面为二维码的参数值
	Ticket   string // 二维码的ticket，可用来换取二维码图片
}
