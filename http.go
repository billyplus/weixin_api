package weixin_api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func HttpGet[T any](url string) (*T, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error().Err(err).Msg("[reloadGameConfig]新建http请求失败")
		return nil, err
	}

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Error().Err(err).Str("url", url).Msg("[reloadGameConfig]发送http请求失败")
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Error().Err(err).Msg("[reloadGameConfig]无法访问Admin服务器")
		return nil, errors.New(res.Status)
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error().Err(err).Msg("[reloadGameConfig]读取http回包失败")
		return nil, errors.Wrap(err, "无法读取回包")
	}

	var resp T
	if err = json.Unmarshal(data, &resp); err != nil {
		log.Error().Bytes("data", data).Msg("[reloadGameConfig]failed parse body")
		return nil, errors.Wrap(err, "无法解析回包")
	}

	return &resp, nil
}

func PostJSON[T any](url string, body interface{}) (*T, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, errors.Wrap(err, "json.Marshal:")
	}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		log.Error().Err(err).Msg("[reloadGameConfig]新建http请求失败")
		return nil, err
	}

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Error().Err(err).Str("url", url).Msg("[reloadGameConfig]发送http请求失败")
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Error().Err(err).Msg("[reloadGameConfig]无法访问Admin服务器")
		return nil, errors.New(res.Status)
	}
	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error().Err(err).Msg("[reloadGameConfig]读取http回包失败")
		return nil, errors.Wrap(err, "无法读取回包")
	}

	var resp T
	if err = json.Unmarshal(data, &resp); err != nil {
		log.Error().Bytes("data", data).Msg("[reloadGameConfig]failed parse body")
		return nil, errors.Wrap(err, "无法解析回包")
	}

	return &resp, nil
}
