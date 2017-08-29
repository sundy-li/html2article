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

	maxDens        float64
	sn             float64
	swn            float64
	title          string
	accurateTitle  string
	atTitleMatched bool

	option *Option
}

type Option struct {
	RemoveNoise   bool // remove noise node
	AccurateTitle bool // find the accurate title node
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
	ext = &extractor{data: make(map[*Info]*html.Node), doc: doc, option: DEFAULT_OPTION}
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
	DEFAULT_OPTION = &Option{
		RemoveNoise: true,
	}
)

func (ec *extractor) SetOption(option *Option) {
	ec.option = option
}

func (ec *extractor) ToArticle() (article *Article, err error) {
	body := find(ec.doc, isTag(atom.Body))
	if body == nil {
		body = ec.doc
	}

	titleNode := find(ec.doc, isTag(atom.Title))
	if titleNode != nil {
		ec.title = getText(titleNode)
	}

	ec.getSn()
	ec.getInfo(body)
	if ec.accurateTitle != "" {
		ec.atTitleMatched = true
	}

	node, err := ec.getBestMatch()
	if err != nil {
		return
	}
	if node == nil {
		err = ERROR_NOTFOUND
		return
	}
	if ec.option.RemoveNoise {
		ec.denoise(node)
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
	article.Publishtime = getPublishTime(node)
	//find title
	article.Title = ec.title
	if ec.option.AccurateTitle && ec.atTitleMatched {
		article.Title = ec.accurateTitle
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
		//get the title
		str := strings.TrimSpace(info.Data)
		if ec.option.AccurateTitle && !ec.atTitleMatched && strings.Contains(ec.title, str) {
			switch node.Parent.DataAtom {
			case atom.Meta, atom.Script, atom.Head:
			default:
				if size := len(str); size > len(ec.accurateTitle) {
					ec.accurateTitle = str
				}
			}
		}
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
			info.ImageCount += cInfo.ImageCount
			info.InputCount += cInfo.InputCount
		}

		info.TagCount++

		switch node.DataAtom {
		case atom.A:
			info.LinkTagCount++
			if node.Parent.DataAtom != atom.P {
				info.LinkTextCount += len(info.Data)
			}
		case atom.P:
			info.Pcount++
		case atom.Img, atom.Image:
			info.ImageCount++
		case atom.Input, atom.Textarea, atom.Button:
			info.InputCount++
		}

		if isContentNode(node) {
			ec.addNode(node, info)
		}
		return
	}
	return
}

//正文去噪
//dmin 最小文本密度为正文最大文本密度的3/10
//去噪即删掉正文中文本密度小于dmin的非文本节点
//只清洗前后三个节点
func (ec *extractor) denoise(node *html.Node) {
	dmin := ec.maxDens * 0.3
	var i = 0
	for n := node.FirstChild; n != nil && i < 3; n = n.NextSibling {
		if isNoisingNode(n) {
			info := ec.getInfo(n)
			if info.Density < dmin {
				next := n.NextSibling
				n.Parent.RemoveChild(n)
				n.NextSibling = next
			}
		}
		i++
	}

	i = 0
	for n := node.LastChild; n != nil && i < 3; n = n.PrevSibling {
		if isNoisingNode(n) {
			info := ec.getInfo(n)
			if info.Density < dmin {
				pre := n.PrevSibling
				n.Parent.RemoveChild(n)
				n.PrevSibling = pre
			}
		}
		i++
	}
}

func (ec *extractor) addNode(node *html.Node, info *Info) {
	info.node = node
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
		if kinfo.score >= maxScore {
			maxScore = kinfo.score
			node = v
		}
		if kinfo.Density > ec.maxDens {
			ec.maxDens = kinfo.Density
		}
	}
	if node == nil {
		err = ERROR_NOTFOUND
	}
	return
}

func getPublishTime(node *html.Node) (ts int64) {
	pnode := node.Parent
	for i := 0; i < 6 && pnode != nil; i++ {
		h, _ := getHtml(pnode)
		ts = getTime(h)
		if ts > 0 {
			break
		}
		pnode = pnode.Parent
	}
	if ts == 0 {
		h, _ := getHtml(node)
		ts = getTime(h)
	}
	return
}

func getAccurateTitle(title string, node *html.Node) (ts int64) {
	pnode := node.Parent
	for i := 0; i < 3 && pnode != nil; i++ {
		for tmpNode := pnode.PrevSibling; tmpNode != nil; tmpNode = tmpNode.PrevSibling {
		}
	}
	if ts == 0 {
		h, _ := getHtml(node)
		ts = getTime(h)
	}
	return
}
