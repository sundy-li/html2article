package html2article

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func DefCode(header http.Header, html string) string {
	contentType := strings.ToLower(header.Get("Content-Type"))
	if strings.Contains(contentType, "charset") {
		return getCharset(contentType)
	}
	// 通过html判断
	index := strings.Index(html, `charset=`)
	if index < 0 || index+20 > len(html) {
		return "utf-8"
	}

	var str = html[index : index+20]
	return getCharset(str)
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
