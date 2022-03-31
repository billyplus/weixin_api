package weixin_api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleMessage(t *testing.T) {
	e := New(&WeiXinApiConfig{})
	body := []byte("\u003cxml\u003e\u003cToUserName\u003e\u003c![CDATA[gh_b72dace7720c]]\u003e\u003c/ToUserName\u003e\n\u003cFromUserName\u003e\u003c![CDATA[o5Iu36FyFfQerLL8fcImYZ40Sstk]]\u003e\u003c/FromUserName\u003e\n\u003cCreateTime\u003e1648698605\u003c/CreateTime\u003e\n\u003cMsgType\u003e\u003c![CDATA[event]]\u003e\u003c/MsgType\u003e\n\u003cEvent\u003e\u003c![CDATA[unsubscribe]]\u003e\u003c/Event\u003e\n\u003cEventKey\u003e\u003c![CDATA[]]\u003e\u003c/EventKey\u003e\n\u003c/xml\u003e")
	e.HandleMessage(body)
	assert.True(t, false)
}
