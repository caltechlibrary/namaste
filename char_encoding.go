package namaste

import (
	"strings"
)

var (
	stepOneEncoding = map[string]string{
		"^": "^5e",
	}
	stepTwoEncoding = map[string]string{
		"\"": "^22",
		"*":  "^2A",
		"/":  "^2F",
		":":  "^3A",
		"<":  "^3C",
		">":  "^3E",
		"?":  "^3F",
		"\\": "^5C",
		"|":  "^7C",
	}
)

func charEncode(s string) string {
	//NOTE: we need to replace ^ with ^5e and avoid collisions with other hex values
	// we split the string into an array of substrings then replace each one as as need to.
	p := strings.Split(s, "")
	for i, target := range p {
		if val, ok := stepOneEncoding[target]; ok == true {
			p[i] = val
		}
	}
	s = strings.Join(p, "")
	for target, replacement := range stepTwoEncoding {
		if strings.Contains(s, target) {
			s = strings.Replace(s, target, replacement, -1)
		}
	}
	return s
}

func charDecode(s string) string {
	for replacement, target := range stepTwoEncoding {
		if strings.Contains(s, target) {
			s = strings.Replace(s, target, replacement, -1)
		}
	}
	for replacement, target := range stepOneEncoding {
		if strings.Contains(s, target) {
			s = strings.Replace(s, target, replacement, -1)
		}
	}
	return s
}
