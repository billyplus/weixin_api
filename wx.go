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

type IEngine interface {
	GetAccessToken() (string, error)
}

var _ IEngine = (*Engine)(nil)

type Engine struct {
	// accessToken       string
	// accessExpiredTime time.Time
	appId     string
	appSecret string
	appToken  string
	wxDomain  string
	// accessToken            string
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
	AppId        string
	AppSecret    string
	AppToken     string
	WeiXinDomain string
	Repository   IRepository
	// AccessToken            string
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
	e.appToken = cfg.AppToken
	e.appSecret = cfg.AppSecret
	e.repo = cfg.Repository

	if cfg.WeiXinDomain == "" {
		e.wxDomain = "https://api.weixin.qq.com"
	} else {
		e.wxDomain = "https://" + cfg.WeiXinDomain
	}

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
	e.client = &http.Client{}
	return e
}

func (e *Engine) GetAccessToken() (string, error) {
	tok, expire, err := e.repo.GetAccessToken(context.Background())
	if err != nil {
		return "", errors.WithMessage(err, "repo.GetAccessToken")
	}
	if tok == "" || time.Now().After(expire) {
		// tok?????????tok?????????
		if err = e.GrantAccessToken(); err != nil {
			if tok != "" {
				// ????????????token????????????????????????tok?????????????????????
				return tok, nil
			}
			return "", errors.WithMessage(err, "GrantAccessToken")
		}

		// ????????????access token
		tok, _, err = e.repo.GetAccessToken(context.Background())
		if err != nil {
			return "", errors.WithMessage(err, "repo.GetAccessToken")
		}

	}
	return tok, nil
}

type responseGrantToken struct {
	ErrorMsg
	AccessToken string `json:"access_token"`
	ExpiresIn   int32  `json:"expires_in"`
}

// ????????????????????????Access Token???????????????repository?????????????????????GetAccessToken????????????repository????????????
func (e *Engine) GrantAccessToken() error {
	if err := e.repo.Lock(); err != nil {
		return errors.WithMessage(err, "repo.Lock")
	}
	defer e.repo.UnLock()
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

	// ??????60?????????
	if err = e.repo.UpdateAccessToken(context.Background(), v.AccessToken, time.Now().Add(time.Duration(v.ExpiresIn-60)*time.Second)); err != nil {
		return errors.Wrap(err, "repo.UpdateAccessToken")
	}

	return nil
}

// ??????Access????????????
// func (e *Engine) AccessExpired() bool {
// 	return e.accessExpiredTime.Before(time.Now())
// }

type QRCodeInfo struct {
	ErrorMsg
	Ticket    string `json:"ticket"`
	ExpiresIn int32  `json:"expire_seconds"`
	URL       string `json:"url"`
}

type qrCodeActionInfo[T any] struct {
	Scene map[string]T `json:"scene"`
}

type qrCodeReqBody[T any] struct {
	ExpireSeconds int32               `json:"expire_seconds,omitempty"`
	ActionName    string              `json:"action_name"`
	ActionInfo    qrCodeActionInfo[T] `json:"action_info"`
}

func CreateQRCode(e IEngine, id int32, expireSeconds int32) (*QRCodeInfo, error) {
	return createQRCode(e, "QR_SCENE", "scene_id", id, expireSeconds)
}

func CreateQRCodeByStr(e IEngine, id string, expireSeconds int32) (*QRCodeInfo, error) {
	return createQRCode(e, "QR_STR_SCENE", "scene_str", id, expireSeconds)
}

func CreateLimitQRCode(e IEngine, id int32) (*QRCodeInfo, error) {
	return createQRCode(e, "QR_LIMIT_SCENE", "scene_id", id, 0)
}

func CreateLimitQRCodeByStr(e IEngine, id string) (*QRCodeInfo, error) {
	return createQRCode(e, "QR_LIMIT_STR_SCENE", "scene_str", id, 0)
}

func createQRCode[IdType any](e IEngine, actionName, idKey string, id IdType, expireSeconds int32) (*QRCodeInfo, error) {
	tok, err := e.GetAccessToken()
	if err != nil {
		return nil, errors.WithMessage(err, "GetAccessToken:")
	}
	req := qrCodeReqBody[IdType]{
		ExpireSeconds: expireSeconds,
		ActionName:    actionName,
		ActionInfo: qrCodeActionInfo[IdType]{
			Scene: map[string]IdType{
				idKey: id,
			},
		},
	}

	info, err := PostJSON[QRCodeInfo](`https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token=`+tok, &req)
	if err != nil {
		return nil, errors.WithMessage(err, "PostJSON:")
	}

	if info.ErrCode > 0 {
		return nil, errors.WithStack(info)
	}

	return info, nil
}

type SessionInfo struct {
	ErrorMsg
	OpenId     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
}

func (e *Engine) Code2Session(jscode string) (*SessionInfo, error) {
	// url: https://api.weixin.qq.com/sns/jscode2session?appid=APPID&secret=SECRET&js_code=JSCODE&grant_type=authorization_code

	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", e.appId, e.appSecret, jscode)
	info, err := HttpGet[SessionInfo](url)
	if err != nil {
		return nil, errors.WithMessage(err, "PostJSON:")
	}

	if info.ErrCode > 0 {
		return nil, errors.WithStack(info)
	}

	return info, nil
}

type reqPhoneNumberBody struct {
	Code string `json:"code"`
}

type PhoneInfo struct {
	PhoneNumber     string
	PurePhoneNumber string
	CountryCode     string
	Watermark       Watermark
}

type respPhoneNumber struct {
	ErrorMsg
	PhoneInfo *PhoneInfo `json:"phone_info"`
}

func (e *Engine) GetPhoneNumber(code string) (*PhoneInfo, error) {
	// url: https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token=ACCESS_TOKEN
	tok, err := e.GetAccessToken()
	if err != nil {
		return nil, errors.WithMessage(err, "GetAccessToken:")
	}
	req := reqPhoneNumberBody{
		Code: code,
	}

	url := fmt.Sprintf("https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token=%s", tok)
	info, err := PostJSON[respPhoneNumber](url, &req)
	if err != nil {
		return nil, errors.WithMessage(err, "PostJSON:")
	}

	if info.ErrCode > 0 {
		return nil, errors.WithStack(info)
	}

	return info.PhoneInfo, nil
}
