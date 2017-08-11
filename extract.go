package html2article

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type extractor struct {
	data   map[*Info]*html.Node
	urlStr string
	doc    *html.Node

	sn  float64
	swn float64
}

func NewFromHtml(htmlStr string) (ext *extractor, err error) {
	return NewFromReader(strings.NewReader(htmlStr))
}

func NewFromReader(reader io.Reader) (ext *extractor, err error) {
	doc, err := html.Parse(reader)
	if err != nil {
		return
	}
	return NewFromNode(doc)
}

func NewFromNode(doc *html.Node) (ext *extractor, err error) {
	ext = &extractor{data: make(map[*Info]*html.Node), doc: doc}
	return
}

func NewFromUrl(urlStr string) (ext *extractor, err error) {
	resp, err := http.Get(urlStr)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	bs, _ := ioutil.ReadAll(resp.Body)
	htmlStr := string(bs)
	htmlStr = DecodeHtml(resp.Header, htmlStr, htmlStr)
	ext, err = NewFromHtml(htmlStr)
	if err != nil {
		return
	}
	ext.urlStr = urlStr
	return
}

var (
	ERROR_NOTFOUND = errors.New("Content not found")
)

func (ec *extractor) ToArticle() (article *Article, err error) {
	body := find(ec.doc, isTag(atom.Body))
	ec.getSn()
	ec.getInfo(body)
	node, err := ec.getBestMatch()
	if err != nil {
		return
	}
	if node == nil {
		err = ERROR_NOTFOUND
		return
	}
	article = &Article{}
	// Get the Content
	article.contentNode = node
	article.Content = getText(node)
	article.Html, err = getHtml(node)
	if err != nil {
		return
	}
	article.Images = getImages(node)
	pnode := node.Parent
	filterNode := func(n *html.Node) bool {
		return n != node
	}
	for i := 0; i < 6 && pnode != nil; i++ {
		article.Publishtime = getTime(getText(pnode, filterNode))
		if article.Publishtime > 0 {
			break
		}
		pnode = pnode.Parent
	}
	if article.Publishtime == 0 {
		article.Publishtime = getTime(article.Content)
	}
	titleNode := find(ec.doc, isTag(atom.Title))
	if titleNode != nil {
		article.Title = getText(titleNode)
	}
	article.Images = getImages(node)
	return
}

func (ec *extractor) getSn() {
	txt := text(ec.doc)
	ec.swn = float64(countStopWords(txt))
	ec.sn = float64(countSn(txt))
}

func (ec *extractor) getInfo(node *html.Node) (info *Info) {
	info = NewInfo()
	if node.Type == html.TextNode {
		info.TextCount = len(node.Data)
		info.LeafList = append(info.LeafList, info.TextCount)
		info.Data = node.Data

		ec.addNode(node, info)
		return
	} else if node.Type == html.ElementNode {
		if isTag(atom.Style)(node) || isTag(atom.Script)(node) {
			return
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			cInfo := ec.getInfo(c)
			info.TextCount += cInfo.TextCount
			info.LinkTextCount += cInfo.LinkTextCount
			info.TagCount += cInfo.TagCount
			info.LinkTagCount += cInfo.LinkTagCount
			info.LeafList = append(info.LeafList, cInfo.LeafList...)
			info.Data += cInfo.Data
			info.Pcount += cInfo.Pcount
		}

		info.TagCount++

		if isTag(atom.A)(node) {
			info.LinkTagCount++
			info.LinkTextCount += len(info.Data)
		} else if isTag(atom.P)(node) {
			info.Pcount++
		}
		if isContentNode(node) {
			ec.addNode(node, info)
		}
		return
	}
	return
}

func (ec *extractor) addNode(node *html.Node, info *Info) {
	info.CalScore(ec.sn, ec.swn)
	ec.data[info] = node
}

func (ec *extractor) getBestMatch() (node *html.Node, err error) {
	if len(ec.data) < 1 {
		err = ERROR_NOTFOUND
		return
	}
	var maxScore float64 = -100
	for kinfo, v := range ec.data {
		// //wechat
		// if cls := attr(v, "id"); cls == "js_content" {
		// 	node = v
		// 	return
		// }
		if kinfo.score >= maxScore {
			maxScore = kinfo.score
			node = v
		}
		// if kinfo.score >= 0 {
		// 	c := attr(v, "class")
		// 	if strings.Contains(c, "article") {
		// 		println("class:", c, kinfo.score, kinfo.Density, kinfo.getAvg(), kinfo.Pcount, kinfo.TextCount, kinfo.LinkTextCount)
		// 	}
		// }
	}
	if node == nil {
		err = ERROR_NOTFOUND
	}
	return
}
