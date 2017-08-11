package html2article

import (
	"net/http"
	"strings"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToArticle(t *testing.T) {
	t.Run("test ToArticle leiphone", func(t *testing.T) {
		assert := assert.New(t)

		testCases := []struct {
			Url         string
			ExpectClass string
		}{
			//雷锋网
			{"https://www.leiphone.com/news/201602/DsiQtR6c1jCu7iwA.html", "lph-article-comView"},
			{"https://www.leiphone.com/news/201708/JQRI6UI8uavpzMwF.html", "lph-article-comView"},
			{"https://www.leiphone.com/news/201708/wlY7YUnEmYHwFFfN.html", "lph-article-comView"},
			{"https://www.leiphone.com/news/201708/DMdo0sSckwJ5FGEv.html", "lph-article-comView"},
			{"https://www.leiphone.com/news/201703/5iXkHxC5eR9VuHpv.html", "lph-article-comView"},
			{"https://www.leiphone.com/news/201708/pIV08b9HKahnoYIM.html", "lph-article-comView"},
			{"https://www.leiphone.com/news/201708/Gs4XTJ1dDPRe328z.html", "lph-article-comView"},
			{"https://www.leiphone.com/news/201707/RRiH46oUsrJSybq0.html", "lph-article-comView"},
			{"https://www.leiphone.com/news/201708/UixD9DKRXaUTts1d.html", "lph-article-comView"},
			{"https://www.leiphone.com/news/201703/OVX5oq3f5jR81wnr.html", "lph-article-comView"},
		}

		for _, testCase := range testCases {
			ext, _ := NewFromUrl(testCase.Url)
			article, err := ext.ToArticle()
			if err != nil {
				t.Fatal(err)
				return
			}
			assert.Nil(err)

			assert.Equal(attr(article.contentNode, "class"), testCase.ExpectClass)
			assert.True(article.Publishtime > 1405732300)
			assert.True(article.Publishtime < 1555732300)
			assert.True(strings.Contains(article.Title, "雷锋网"))
			assert.True(strings.Contains(article.Content, "雷锋网"))
		}

	})

	t.Run("test ToArticle others", func(t *testing.T) {
		assert := assert.New(t)

		testCases := []struct {
			Url         string
			ExpectClass string
		}{
			{"http://cj.sina.com.cn/article/detail/5835524730/241716?column=hkstock&ch=9", "article article_16"},
			{"http://cj.sina.com.cn/article/detail/5617263472/355836?column=stock&ch=9", "article article_16"},
		}

		for _, testCase := range testCases {
			ext, _ := NewFromUrl(testCase.Url)
			article, err := ext.ToArticle()
			if err != nil {
				t.Fatal(err)
				return
			}
			assert.Nil(err)

			assert.Equal(attr(article.contentNode, "class"), testCase.ExpectClass)
			assert.True(article.Publishtime > 1405732300)
			assert.True(article.Publishtime < 1555732300)
		}

	})
}

func BenchmarkToArticle(b *testing.B) {
	urlStr := "https://www.leiphone.com/news/201602/DsiQtR6c1jCu7iwA.html"
	resp, err := http.Get(urlStr)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	for i := 0; i < b.N; i++ {
		ext, err := NewFromReader(resp.Body)
		if err != nil {
			b.Fatal(err.Error())
		}
		ext.ToArticle()
	}
}
