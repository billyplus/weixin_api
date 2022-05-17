package weixin_api

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func (e *Engine) CreateMenu(menu []byte) (*ErrorMsg, error) {
	tok, err := e.GetAccessToken()
	if err != nil {
		return nil, errors.WithMessage(err, "GetAccessToken:")
	}
	// https://api.weixin.qq.com/cgi-bin/user/info?access_token=ACCESS_TOKEN&openid=OPENID&lang=zh_CN
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/menu/create?access_token=%s", tok)

	info, err := PostJSON[ErrorMsg](url, menu)
	if err != nil {
		return nil, errors.WithMessage(err, "HttpGet:")
	}
	return info, nil
}

func (e *Engine) GetCurrentSelfMenuInfo() error {
	tok, err := e.GetAccessToken()
	if err != nil {
		return errors.WithMessage(err, "GetAccessToken:")
	}
	// https://api.weixin.qq.com/cgi-bin/user/info?access_token=ACCESS_TOKEN&openid=OPENID&lang=zh_CN
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/get_current_selfmenu_info?access_token=%s", tok)

	info, err := HttpGetRaw(url)
	if err != nil {
		return errors.WithMessage(err, "HttpGet:")
	}
	log.Debug().Str("menu", string(info)).Msg("GetCurrentSelfMenuInfo")
	return nil
}
