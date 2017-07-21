package html2article

import (
	"net/url"
	"path"
	"strings"
)

type Article struct {
	// Html is content with html tag
	Html        string   `json:"content_html"`
	Content     string   `json:"content"`
	Title       string   `json:"title"`
	Publishtime int64    `json:"publish_time"`
	Images      []string `json:"images"`
}

// ParseImage parse the image src to the absolute path
func (a *Article) ParseImage(urlStr string) {
	_url, err := url.Parse(urlStr)
	if err != nil {
		return
	}
	mp := make(map[string]string)
	for i, _ := range a.Images {
		if strings.Index(a.Images[i], "http") != 0 {
			var newImg string
			if strings.Index(a.Images[i], "//") == 0 {
				newImg = _url.Scheme + ":" + a.Images[i]
			} else if strings.Index(a.Images[i], "/") == 0 {
				newImg = _url.Scheme + "://" + _url.Host + a.Images[i]
			} else {
				newImg = _url.Scheme + "://" + _url.Host + path.Join(_url.Path, "../", a.Images[i])
			}
			mp[a.Images[i]] = newImg
			a.Images[i] = newImg
		}
	}
	for k, v := range mp {
		a.Html = strings.Replace(a.Html, k, v, -1)
	}
}
