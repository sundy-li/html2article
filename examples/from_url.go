package main

import (
	"github.com/sundy-li/html2article"
)

func main() {
	article, err := html2article.FromUrl("http://edition.cnn.com/travel/article/airlines-cabin-waste/index.html")
	if err != nil {
		panic(err)
	}
	println("article title is =>", article.Title)
	println("article publishtime is =>", article.Publishtime)
	println("article content is =>", article.Content)
}
