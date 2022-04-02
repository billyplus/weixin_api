package weixin_api

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleMessage(t *testing.T) {
	e := New(&WeiXinApiConfig{})
	body := []byte(`<xml>
	<ToUserName><![CDATA[toUser]]></ToUserName>
	<FromUserName><![CDATA[FromUser]]></FromUserName>
	<CreateTime>123456789</CreateTime>
	<MsgType><![CDATA[event]]></MsgType>
	<Event><![CDATA[subscribe]]></Event>
	<EventKey><![CDATA[qrscene_123123]]></EventKey>
	<Ticket><![CDATA[TICKET]]></Ticket>
  </xml>`)
	err := e.HandleMessage(context.TODO(), body)
	assert.Nil(t, err, "should not return error")
	assert.True(t, false)
}
