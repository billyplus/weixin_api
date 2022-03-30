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
