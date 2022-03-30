package weixinapi

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
	appId     string
	appSecret string
	wxDomain  string
	repo      IRepository
	client    *http.Client
}

type WeiXinApiConfig struct {
	AppId        string
	AppSecret    string
	WeiXinDomain string
	AccessToken  string
}

func New(cfg *WeiXinApiConfig) *Engine {
	e := &Engine{}
	e.appId = cfg.AppId
	e.appSecret = cfg.AppSecret
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
