package html2article

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefCode(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		in       string
		expected string
	}{
		// http://xiaozhenkai.blog.51cto.com/1029756/1949841
		{`<script id="allmobilize" charset="utf-8" src="http://a.yunshipei.com/ac0ecd4968d76dbea889a1a0f28e176f/allmobilize.min.js"></script><meta http-equiv="Cache-Control" content="no-siteapp" /><link rel="alternate" media="handheld" href="#" /><meta http-equiv="Content-Type" content="text/html;charset=gb2312">`, "gbk"},
		{`<meta charset="utf-8">`, "utf-8"},
	}

	for _, tt := range testCases {
		actual := DefCode(map[string][]string{}, tt.in)
		assert.Equal(tt.expected, actual)
	}
}
