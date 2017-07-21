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
	data map[*Info]*html.Node
}

func NewExtractor() *extractor {
	return &extractor{data: make(map[*Info]*html.Node)}
}

var (
	ERROR_NOTFOUND = errors.New("Content not found")
)

// FromUrl do parse the urlStr html to  Article
func FromUrl(urlStr string) (article *Article, err error) {
	req, _ := http.NewRequest("GET", urlStr, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.98 Safari/537.36")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	bs, _ := ioutil.ReadAll(resp.Body)

	htmlStr := string(bs)
	htmlStr = DecodeHtml(resp.Header, htmlStr, htmlStr)
	article, err = FromReader(strings.NewReader(htmlStr))
	return
}

// FromHtml do parse the htmlStr to  Article
func FromHtml(htmlStr string) (article *Article, err error) {
	return FromReader(strings.NewReader(htmlStr))
}

// FromNode
func FromNode(node *html.Node) (article *Article, err error) {
	return extract(node)
}

// From Reader
func FromReader(reader io.Reader) (article *Article, err error) {
	doc, err := html.Parse(reader)
	if err != nil {
		return
	}
	return extract(doc)
}

func extract(doc *html.Node) (article *Article, err error) {
	ec := NewExtractor()
	body := find(doc, isTag(atom.Body))
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
	titleNode := find(doc, isTag(atom.Title))
	if titleNode != nil {
		article.Title = getText(titleNode)
	}
	article.Images = getImages(node)
	return
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
			info.DensitySum += cInfo.Density
			info.Data += cInfo.Data
		}

		info.TagCount++

		if isTag(atom.A)(node) {
			info.LinkTagCount++
		} else if isTag(atom.P)(node) {
			info.Pcount++
		}
		var9 := info.TextCount - info.LinkTextCount
		var10 := info.TagCount - info.LinkTagCount
		if var9*var10 != 0 {
			info.Density = (float64(var9) / float64(var10))
		}
		if isContentNode(node) {
			ec.addNode(node, info)
		}
		return
	}
	return
}

func (ec *extractor) addNode(node *html.Node, info *Info) {
	info.CalScore()
	ec.data[info] = node
	// cls := attr(node, "class")
	// if cls == "bk3left" {
	// 	fmt.Printf("%v\n", info)
	// }
}

func (ec *extractor) getBestMatch() (node *html.Node, err error) {
	if len(ec.data) < 1 {
		err = ERROR_NOTFOUND
		return
	}

	var maxScore float64 = -100
	for kinfo, v := range ec.data {
		//wechat
		if cls := attr(v, "id"); cls == "js_content" {
			node = v
			return
		}

		//如果不含有标点符号,那么不算入正文
		if !strings.Contains(kinfo.Data, "，") && !strings.Contains(kinfo.Data, ",") {
			continue
		}
		if kinfo.score >= maxScore {
			maxScore = kinfo.score
			node = v
		}
	}
	if node == nil {
		err = ERROR_NOTFOUND
	}
	return
}
