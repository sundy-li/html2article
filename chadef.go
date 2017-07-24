package html2article

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

var charsetRe = regexp.MustCompile(`<meta\s.*charset=(.*)>`)

func DefCode(header http.Header, html string) string {
	contentType := strings.ToLower(header.Get("Content-Type"))
	if strings.Contains(contentType, "charset") {
		return getCharset(contentType)
	}
	// 通过html判断
	matches := charsetRe.FindStringSubmatch(html)
	if len(matches) > 1 {
		return getCharset(strings.Trim(matches[1], `"`))
	}

	return "utf-8"
}

func getCharset(contentType string) string {
	contentType = strings.ToLower(contentType)
	if strings.Contains(contentType, "utf-8") {
		return "utf-8"
	}
	if strings.Contains(contentType, "gbk") {
		return "gbk"
	}
	if strings.Contains(contentType, "gb-18030") {
		return "gb-18030"
	}
	if strings.Contains(contentType, "gb") {
		return "gbk"
	}
	return "utf-8"
}

func DecodeHtml(header http.Header, word, src string) (dst string) {
	typ := DefCode(header, src)
	var encoder encoding.Encoding
	switch typ {
	case "utf-8":
		return word
	case "gbk":
		encoder = simplifiedchinese.GBK
	case "gb-18030":
		encoder = simplifiedchinese.GB18030
	default:
		encoder = simplifiedchinese.GBK
	}
	data, err := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(word)), encoder.NewDecoder()))
	if err == nil {
		dst = string(data)
	}
	dst = string(data)
	return
}
