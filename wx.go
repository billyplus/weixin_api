package weixin_api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type Engine struct {
	// accessToken       string
	// accessExpiredTime time.Time
	appId                  string
	appSecret              string
	wxDomain               string
	repo                   IRepository
	client                 *http.Client
	handleTextMessage      func(m *TextMessage) error
	handleImageMessage     func(m *ImageMessage) error
	handleVoiceMessage     func(m *VoiceMessage) error
	handleVideoMessage     func(m *VideoMessage) error
	handleLocationMessage  func(m *LocationMessage) error
	handleLinkMessage      func(m *LinkMessage) error
	handleClickEvent       func(m *ClickEvent) error
	handleLocationEvent    func(m *LocationEvent) error
	handleViewEvent        func(m *ViewEvent) error
	handleScanEvent        func(m *ScanEvent) error
	handleSubscribeEvent   func(m *SubscribeEvent) error
	handleUnsubscribeEvent func(m *UnsubscribeEvent) error
}

type WeiXinApiConfig struct {
	AppId                  string
	AppSecret              string
	WeiXinDomain           string
	AccessToken            string
	HandleTextMessage      func(m *TextMessage) error
	HandleImageMessage     func(m *ImageMessage) error
	HandleVoiceMessage     func(m *VoiceMessage) error
	HandleVideoMessage     func(m *VideoMessage) error
	HandleLocationMessage  func(m *LocationMessage) error
	HandleLinkMessage      func(m *LinkMessage) error
	HandleClickEvent       func(m *ClickEvent) error
	HandleLocationEvent    func(m *LocationEvent) error
	HandleViewEvent        func(m *ViewEvent) error
	HandleScanEvent        func(m *ScanEvent) error
	HandleSubscribeEvent   func(m *SubscribeEvent) error
	HandleUnsubscribeEvent func(m *UnsubscribeEvent) error
}

func New(cfg *WeiXinApiConfig) *Engine {
	e := &Engine{}
	e.appId = cfg.AppId

	e.handleTextMessage = cfg.HandleTextMessage
	e.handleImageMessage = cfg.HandleImageMessage
	e.handleVoiceMessage = cfg.HandleVoiceMessage
	e.handleVideoMessage = cfg.HandleVideoMessage
	e.handleLocationMessage = cfg.HandleLocationMessage
	e.handleLinkMessage = cfg.HandleLinkMessage
	e.handleClickEvent = cfg.HandleClickEvent
	e.handleLocationEvent = cfg.HandleLocationEvent
	e.handleViewEvent = cfg.HandleViewEvent
	e.handleScanEvent = cfg.HandleScanEvent
	e.handleSubscribeEvent = cfg.HandleSubscribeEvent
	e.handleUnsubscribeEvent = cfg.HandleUnsubscribeEvent
	if cfg.WeiXinDomain == "" {
		e.wxDomain = "https://api.weixin.qq.com"
	} else {
		e.wxDomain = "https://" + cfg.WeiXinDomain
	}
	// e.accessToken = cfg.AccessToken
	e.client = &http.Client{}
	return e
}

func (e *Engine) GetAccessToken() (string, error) {
	tok, err := e.repo.GetAccessToken(context.TODO())
	if err != nil {
		return "", errors.WithMessage(err, "repo.GetAccessToken")
	}
	return tok, nil
}

type responseGrantToken struct {
	ErrorMsg
	AccessToken string `json:"access_token"`
	ExpiresIn   int32  `json:"expires_in"`
}

// 从微信服务器获取Access Token
func (e *Engine) GrantAccessToken() error {
	// https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=APPID&secret=APPSECRET
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", e.appId, e.appSecret)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return errors.Wrap(err, "http.NewRequest")
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "DoRequest")
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "ReadBody")
	}

	var v responseGrantToken
	if err = json.Unmarshal(data, &v); err != nil {
		return errors.Wrap(err, "UnmarshalBody")
	}

	if v.ErrCode > 0 {
		return errors.WithStack(&v)
	}

	// 提前60秒更新
	if err = e.repo.UpdateAccessToken(context.TODO(), v.AccessToken, time.Now().Add(time.Duration(v.ExpiresIn-60)*time.Second)); err != nil {
		return errors.Wrap(err, "repo.UpdateAccessToken")
	}

	return nil
}

// 检查Access是否过期
// func (e *Engine) AccessExpired() bool {
// 	return e.accessExpiredTime.Before(time.Now())
// }
