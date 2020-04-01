package lodago

import (
	"regexp"
	"strings"
)

// 这里是github.com/gookit/goutil/arrutil库内的实现，感谢作者。

var (
	toSnakeReg  = regexp.MustCompile("[A-Z][a-z]")
	toCamelRegs = map[string]*regexp.Regexp{
		" ": regexp.MustCompile(" +[a-zA-Z]"),
		"-": regexp.MustCompile("-+[a-zA-Z]"),
		"_": regexp.MustCompile("_+[a-zA-Z]"),
	}
)

// LowerFirst ABC -> aBC
func LowerFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	f := s[0]
	if f >= 'A' && f <= 'Z' {
		return strings.ToLower(string(f)) + s[1:]
	}
	return s
}

// UpperFirst abc -> Abc
func UpperFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	f := s[0]
	if f >= 'a' && f <= 'z' {
		return strings.ToUpper(string(f)) + s[1:]
	}
	return s
}

// CamelCase "group_id" -> "groupId"
func CamelCase(s string, sep ...string) string {
	sepChar := "_"
	if len(sep) > 0 {
		sepChar = sep[0]
	}
	if !strings.Contains(s, sepChar) {
		return s
	}
	rgx, ok := toCamelRegs[sepChar]
	if !ok {
		rgx = regexp.MustCompile(regexp.QuoteMeta(sepChar) + "+[a-zA-Z]")
	}
	return rgx.ReplaceAllStringFunc(s, func(s string) string {
		s = strings.TrimLeft(s, sepChar)
		return UpperFirst(s)
	})
}

// SnakeCase "GroupId" -> "group_id"
func SnakeCase(s string, sep ...string) string {
	sepChar := "_"
	if len(sep) > 0 {
		sepChar = sep[0]
	}
	newStr := toSnakeReg.ReplaceAllStringFunc(s, func(s string) string {
		return sepChar + LowerFirst(s)
	})
	return strings.TrimLeft(newStr, sepChar)
}
