// 用户管理
package weixin_api

import (
	"fmt"

	"github.com/pkg/errors"
)

type UserInfo struct {
}

func (e *Engine) GetUserInfo() (*UserInfo, error) {
	tok, err := e.GetAccessToken()
	if err != nil {
		return nil, errors.WithMessage(err, "GetAccessToken:")
	}
	// https://api.weixin.qq.com/cgi-bin/user/info?access_token=ACCESS_TOKEN&openid=OPENID&lang=zh_CN
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/user/info?access_token=%s&openid=OPENID&lang=zh_CN", tok)

	info, err := HttpGet[UserInfo](url)
	if err != nil {
		return nil, errors.WithMessage(err, "HttpGet:")
	}
	return info, nil
}
