package ctrls

import (
	"fmt"
	"html/template"
	"strconv"
	"strings"
)

func js(in string) (out template.HTML) {
	o := fmt.Sprintf("<script src='/%s/js/%s'></script>", staticDir, in)
	return template.HTML(o)
}

func css(in string) (out template.HTML) {
	o := fmt.Sprintf("<link href='/%s/css/%s'>", staticDir, in)
	return template.HTML(o)
}

func kb(in int64) int64 {
	return in / 1024
}

func isImg(in string) bool {
	return strings.Contains(in, "image")
}

func isActive(in, needle string) string {
	if strings.Contains(in, needle) {
		return "active"
	}
	return ""
}

func CutLongText(s string, n int) string {
	if l := len(s); l <= n {
		return s
	}
	return s[:n] + "..."
}

func urlContain(reqUrl, url string) bool {
	return strings.Contains(reqUrl, url)
}

func floatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 0, 64)
}
