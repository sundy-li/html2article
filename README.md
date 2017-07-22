## 基于文本密度的html2article实现[golang] 

## Install
	go get -u -v github.com/sundy-li/html2article


## Performance
  avg 3.2ms per article, accuracy >= 98% (对比其他开源实现,可能是目前最快的html2article实现,我们测试的数据集约3kw来自于微信公众号,各大类中文科技媒体历史文章,目前能达到98%以上准确率)


## Examples
参考examples
[from_url.go][1]

	
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




## Algorithm
- [参考论文][2]
- [Java实现][3]


[1]: https://github.com/sundy-li/html2article/blob/master/examples/from_url.go
[2]: http://www.doc88.com/p-7714009813182.html
[3]: https://github.com/CrawlScript/WebCollector
 
