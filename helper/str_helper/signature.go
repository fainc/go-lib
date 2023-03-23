package str_helper

import (
	"fmt"
	"sort"
	"strings"
)

type SignatureStrKeyOptions struct {
	Key   string
	Value string
}

// SignatureStr 常用签名处理方法 map转签名串（字典升序），适用一层的map数据
func SignatureStr(m map[string]string, keyOptions *SignatureStrKeyOptions, exclude []string) (str string) {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	// 按字典升序排列
	sort.Strings(keys)
	str = ""
	for _, k := range keys {
		if _, ok := m[k]; ok {
			if exclude == nil || !StrInSlice(k, exclude) {
				str = fmt.Sprintf("%s&%s=%s", str, k, m[k])
			}
		}
	}
	if keyOptions != nil && keyOptions.Key != "" && keyOptions.Value != "" {
		str = fmt.Sprintf("%s&%s=%s", str, keyOptions.Key, keyOptions.Value)
	}
	str = strings.TrimPrefix(str, "&")
	return
}
