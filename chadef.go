package html2article

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/endeveit/enca"
	"github.com/qiniu/iconv"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

var charsetReHeader = regexp.MustCompile(`.*charset=(.*)`)
var charsetRe = regexp.MustCompile(`<meta\s.*charset=(.*)>`)

func DefCode(header http.Header, html string) string {
	contentType := strings.ToLower(header.Get("Content-Type"))
	// 通过header判断
	matches := charsetReHeader.FindStringSubmatch(contentType)
	if len(matches) > 1 {
		return getCharset(strings.Trim(matches[1], `"`))
	}
	// 通过html判断
	matches = charsetRe.FindStringSubmatch(html)
	if len(matches) > 1 {
		return getCharset(strings.Trim(matches[1], `"`))
	}

	analyzer, err := enca.New("zh")
	if err == nil {
		encoding, err := analyzer.FromString(html, enca.NAME_STYLE_ICONV)
		defer analyzer.Free()
		if err == nil {
			return encoding
		}
	}

	return "utf-8"
}

func getCharset(contentType string) string {
	return strings.ToLower(contentType)
}

func DecodeHtml(header http.Header, word, src string) (dst string) {
	typ := DefCode(header, src)
	if typ == "gb2312" {
		typ = "gbk"
	} else if typ == "utf-8" {
		dst = src
		return
	}
	cd, err := iconv.Open("utf-8", typ) // convert typ to utf8
	if err != nil {
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
	defer cd.Close()
	data, err := ioutil.ReadAll(iconv.NewReader(cd, bytes.NewReader([]byte(word)), len(src)))
	dst = string(data)
	return dst
}
