//COPYRIGHT https://github.com/golang/tools/blob/master/cmd/html2article/conv.go
package html2article

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type selector func(*html.Node) bool
type Style string

var (
	newlineRun = regexp.MustCompile(`\s+`)

	timeRegex = []*regexp.Regexp{
		regexp.MustCompile(`([\d]{4})-([\d]{1,2})-([\d]{1,2})\s*([\d]{1,2}:[\d]{1,2})?`),
		regexp.MustCompile(`([\d]{4}).([\d]{1,2}).([\d]{1,2})\s*([\d]{1,2}:[\d]{1,2})?`),
		regexp.MustCompile(`([\d]{4})/([\d]{1,2})/([\d]{1,2})\s*([\d]{1,2}:[\d]{1,2})?`),
		regexp.MustCompile(`([\d]{4})年([\d]{1,2})月([\d]{1,2})日\s*([\d]{1,2}:[\d]{1,2})?`),
	}
)

func countStopWords(str string) int {
	var c = 0
	for k, _ := range stopwords {
		c += strings.Count(str, k)
	}
	if c == 0 {
		c = 1
	}
	return c
}

func countSn(str string) int {
	sn := strings.Count(str, ".")
	sn += strings.Count(str, "。")
	if sn == 0 {
		sn = 1
	}
	return sn
}

func limitNewlineRuns(s string) string {
	return newlineRun.ReplaceAllString(s, " ")
}

func getTime(str string) int64 {
	for _, t := range timeRegex {
		ts := t.FindStringSubmatch(str)
		if len(ts) < 4 {
			continue
		}
		var h string = "00:00"
		if len(ts) > 4 && ts[4] != "" {
			h = ts[4]
		}

		year, _ := strconv.Atoi(ts[1])
		month, _ := strconv.Atoi(ts[2])
		day, _ := strconv.Atoi(ts[3])

		timeAt := strings.Split(h, ":")
		hour, _ := strconv.Atoi(timeAt[0])
		var minute int
		if len(timeAt) > 1 {
			minute, _ = strconv.Atoi(timeAt[1])
		}

		v := fmt.Sprintf("%04d%02d%02d %02d:%02d", year, month, day, hour, minute)
		tm, err := time.Parse("20060102 15:04", v)
		if err == nil {
			return tm.Unix()
		}
	}
	return 0
}

// get Text and transform the charset
func getText(n *html.Node, filter ...selector) string {
	return limitNewlineRuns(strings.TrimSpace(text(n, filter...)))
}

func text(n *html.Node, filter ...selector) string {
	var buf bytes.Buffer
	walk(n, func(n *html.Node) bool {
		if n == nil {
			return false
		}
		switch n.Type {
		case html.TextNode:
			buf.WriteString(n.Data)
			return false
		case html.ElementNode:
			// no-op
		default:
			return true
		}
		if isTag(atom.Style)(n) || isTag(atom.Script)(n) || isTag(atom.Image)(n) || isTag(atom.Img)(n) || isTag(atom.Textarea)(n) || isTag(atom.Input)(n) || isTag(atom.Noscript)(n) {
			return false
		}
		buf.WriteString(childText(n, filter...))
		return false
	})
	return buf.String()
}

func childText(node *html.Node, filter ...selector) string {
	var buf bytes.Buffer
	for n := node.FirstChild; n != nil; n = n.NextSibling {
		flag := true
		for _, f := range filter {
			flag = flag && f(node)
		}
		if flag {
			fmt.Fprint(&buf, text(n, filter...))
		}
	}
	return buf.String()
}

func getHtml(n *html.Node) (str string, err error) {
	var buf bytes.Buffer
	err = html.Render(&buf, n)
	str = buf.String()
	return
}

func getImages(node *html.Node) []string {
	res := []string{}
	mp := make(map[string]bool)
	walk(node, func(n *html.Node) bool {
		if isTag(atom.Img)(n) {
			if width, err := strconv.Atoi(attr(n, "width")); err == nil {
				if width != 0 && width < 30 {
					return false
				}
			}

			if height, err := strconv.Atoi(attr(n, "height")); err == nil {
				if height != 0 && height < 30 {
					return false
				}
			}

			// 不抓取默认不展示图片
			if display := attr(n.Parent, "style"); len(display) > 0 && strings.Contains(display, "display: none") {
				return false
			}

			if display := attr(n, "style"); len(display) > 0 && strings.Contains(display, "display: none") {
				return false
			}

			src := attr(n, "data-original")

			if len(src) == 0 {
				src = attr(n, "src")
			} else {
				setAttr(n, "src", src)
				removeAttr(n, "data-src")
			}

			if len(src) == 0 {
				src = attr(n, "data-src")
			}

			excludeStrs := []string{
				"w16_h16.png",
				"logo.png",
				"icon.png",
			}

			if len(src) > 0 {
				for _, exc := range excludeStrs {
					if strings.Contains(src, exc) {
						return false
					}
				}
			}

			if _, ok := mp[src]; !ok && len(src) > 0 {
				mp[src] = true
				res = append(res, src)
			}
			return false
		} else {
			for c := node.FirstChild; c != nil; c = c.NextSibling {
				res = append(res, getImages(c)...)
			}
			return false
		}
	})
	return res
}

func isContentNode(n *html.Node) bool {
	return n.DataAtom == atom.Div || n.DataAtom == atom.Section || n.DataAtom == atom.Table || n.DataAtom == atom.Tr || n.DataAtom == atom.Td || n.DataAtom == atom.Tbody || n.DataAtom == atom.Article
}

func isNoisingNode(n *html.Node) bool {
	// isnoreturn n.DataAtom != atom.A && n.DataAtom != atom.Image && n.DataAtom != atom.P
	return n.DataAtom == atom.Div || n.DataAtom == atom.Script || n.DataAtom == atom.Section
}

func isTag(a atom.Atom) selector {
	return func(n *html.Node) bool {
		return n.DataAtom == a
	}
}

// func hasContent(str string) selector {
// 	return func(n *html.Node) bool {
// 		return n.Data
// 	}
// }

func alwaysTrue() selector {
	return func(n *html.Node) bool {
		return true
	}
}

func hasAttr(key, val string) selector {
	return func(n *html.Node) bool {
		for _, a := range n.Attr {
			if a.Key == key && a.Val == val {
				return true
			}
		}
		return false
	}
}

func attr(node *html.Node, key string) (value string) {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func findAll(node *html.Node, fn selector) (nodes []*html.Node) {
	walk(node, func(n *html.Node) bool {
		if fn(n) {
			nodes = append(nodes, n)
		}
		return true
	})
	return
}

func find(n *html.Node, fn selector) *html.Node {
	var result *html.Node
	walk(n, func(n *html.Node) bool {
		if result != nil {
			return false
		}
		if fn(n) {
			result = n
			return false
		}
		return true
	})
	return result
}

func walk(n *html.Node, fn selector) {
	if fn(n) {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c, fn)
		}
	}
}

func removeAttr(n *html.Node, attrName string) {
	for i, a := range n.Attr {
		if a.Key == attrName {
			n.Attr[i], n.Attr[len(n.Attr)-1], n.Attr =
				n.Attr[len(n.Attr)-1], html.Attribute{}, n.Attr[:len(n.Attr)-1]
			return
		}
	}
}

func setAttr(n *html.Node, attrName, value string) {
	for i, a := range n.Attr {
		if a.Key == attrName {
			n.Attr[i].Val = value
			return
		}
	}
	if len(n.Attr) == 0 {
		n.Attr = make([]html.Attribute, 1)
	}
	n.Attr = append(n.Attr, html.Attribute{
		Key: attrName,
		Val: value,
	})
}
