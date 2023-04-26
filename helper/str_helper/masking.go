package str_helper

import (
	"regexp"
)

type MaskParams struct {
	Key []string // 关键词遮掩
	Pos []int    // 字符位置遮掩，pos长度决定字符位置遮掩规则，长度 1:[N] 遮掩前N个字符;长度2:[N,M] 遮掩N-M区间(含)的字符;长度>2:[X,Y,Z],自定义字符位置遮掩
}

func MaskMobile(m string, rule int) string {
	var r []int
	switch rule {
	case 1: // 156******01
		r = []int{3, 8}
	case 2: // 156********
		r = []int{3, 10}
	default: // 156****7601
		r = []int{3, 6}
	}
	return Mask(m, &MaskParams{Pos: r})
}
func MaskIdentity(m string, rule int) string {
	var r []int
	switch rule {
	case 1: // 15************6428
		r = []int{2, 13}
	case 2: // 15**************28
		r = []int{2, 15}
	case 3: // **************6428
		r = []int{0, 13}
	case 4: // ****************28 Strong
		r = []int{0, 15}
	case 5: // 35**************** Strong
		r = []int{2, 17}
	default: // 152530********6428 Weak
		r = []int{6, 13}
	}
	return Mask(m, &MaskParams{Pos: r})
}
func MaskAddress(add string, _ int) string {
	var numbersRegExp = regexp.MustCompile("[0-9]+")
	return numbersRegExp.ReplaceAllString(add, "*")
}

// Mask 字符串遮掩操作（不支持中文）
func Mask(str string, p *MaskParams) string {
	strLen := len(str)
	if p != nil && p.Pos != nil {
		if p.Pos[0] >= strLen {
			p.Pos[0] = strLen - 1
		}
		if p.Pos[0] < 0 {
			p.Pos[0] = 0
		}
		if len(p.Pos) == 1 {
			str = maskMarker(p.Pos[0]) + str[p.Pos[0]+1:]
		}
		if len(p.Pos) == 2 {
			if p.Pos[1] < p.Pos[0] {
				p.Pos[1] = p.Pos[0]
			}
			if p.Pos[1] >= strLen {
				p.Pos[1] = strLen - 1
			}
			if p.Pos[1] < 0 {
				p.Pos[1] = 0
			}
			str = str[:p.Pos[0]] + maskMarker(p.Pos[1]-p.Pos[0]) + str[p.Pos[1]+1:]
		}
		if len(p.Pos) > 2 {
			ns := []byte(str)
			for _, v := range p.Pos {
				ns[v] = '*'
			}
			str = string(ns)
		}
	}
	return str
}

func maskMarker(n int) (ret string) {
	ret = "*"
	for i := 0; i < n; i++ {
		ret = ret + "*"
	}
	return
}
