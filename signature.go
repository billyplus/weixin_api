package weixin_api

import (
	"crypto/sha1"
	"fmt"
	"sort"
	"strings"
)

// 验证签名是否合法
func ValidateSignature(tok, timestamp, nonce, signature string) bool {
	strs := []string{tok, timestamp, nonce}
	sort.Strings(strs)

	tmpStr := strings.Join(strs, "")
	actual := fmt.Sprintf("%x", sha1.Sum([]byte(tmpStr)))

	return actual == signature
}

// 验证签名是否合法
func (e *Engine) ValidateSignature(timestamp, nonce, signature string) bool {
	return ValidateSignature(e.appToken, timestamp, nonce, signature)
}

type Watermark struct {
	AppId     string `json:"appid"`     // 应用的appid
	Timestamp int64  `json:"timestamp"` // 操作的时间戳
}
