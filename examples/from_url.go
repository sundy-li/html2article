package main

import (
	"github.com/sundy-li/html2article"
)

func main() {
	article, err := html2article.FromUrl("https://www.leiphone.com/news/201602/DsiQtR6c1jCu7iwA.html")
	if err != nil {
		panic(err)
	}
	println("article title is =>", article.Title)
	println("article publishtime is =>", article.Publishtime)
	println("article content is =>", article.Content)
}
