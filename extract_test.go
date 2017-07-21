package html2article

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtract(t *testing.T) {
	t.Run("test extract main", func(t *testing.T) {
		assert := assert.New(t)
		urlStr := "https://www.leiphone.com/news/201602/DsiQtR6c1jCu7iwA.html"

		article, err := FromUrl(urlStr)
		if err != nil {
			t.Fatal(err)
			return
		}
		assert.Nil(err)
		assert.Equal(int64(1455732300), article.Publishtime)
		assert.Equal(3, len(article.Images))
		assert.Equal("“朋友印象”被曝可以看到匿名者真实身份，还能否愉快地八卦？ | 雷锋网", article.Title)
		assert.Contains(article.Content, "春节这几天")
		assert.Contains(article.Html, "<p>")
		assert.Contains(article.Html, "春节这几天")
	})
}

func BenchmarkExtract(b *testing.B) {
	urlStr := "http://tech.qq.com/a/20161205/001738.htm"
	resp, err := http.Get(urlStr)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	bs, _ := ioutil.ReadAll(resp.Body)
	for i := 0; i < b.N; i++ {
		FromReader(bytes.NewReader(bs))
	}
}
