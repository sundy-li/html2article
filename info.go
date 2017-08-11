package html2article

import (
	"math"

	"golang.org/x/net/html"
)

type Info struct {
	TextCount     int
	LinkTextCount int
	TagCount      int
	LinkTagCount  int
	LeafList      []int
	Density       float64
	Pcount        int
	InputCount    int
	ImageCount    int

	Data  string
	score float64

	node *html.Node
}

func NewInfo() *Info {
	return &Info{}
}

func (info *Info) CalScore(sn_sum, swn_sum float64) {
	//avg * ln((cn-lcn)/lcn) * (sn/snm + 1) * (swn/swnm + 1) * abs(ln((cn+1)/(tn+1)) * ln(pn+2) * (tn-in+1)/(tn+1) * (tn-mn+1)/(tn+1) )
	a1 := info.TextCount - info.LinkTextCount
	a2 := info.LinkTextCount
	sn := countSn(info.Data)
	swn := countStopWords(info.Data)

	a3 := math.Abs(math.Log(float64(info.TextCount+1) / float64(info.TagCount+1)))
	a4 := float64(info.TagCount-info.InputCount+1) / float64(info.TagCount+1)
	a5 := float64(info.TagCount-info.ImageCount+1) / float64(info.TagCount+1)

	if a1 == 0 {
		a1 = 1
	}
	if a2 == 0 {
		a2 = 1
	}
	info.Density = math.Log(float64(a1)/float64(a2)) * (float64(sn)/sn_sum + 1) * (float64(swn)/swn_sum + 1) * a3 * a4 * a5
	avg := info.getAvg()

	info.score = math.Log(avg) * float64(info.Density) * math.Log10(float64(info.Pcount+2))
	return
	// if info.score >= 0 {
	// 	c := attr(info.node, "class")
	// 	if strings.Contains(c, "rticle") {
	// 		println("class:", c, info.score, info.Density, avg, math.Log10(float64(info.Pcount+2)), math.Log(float64(a1)/float64(a2)), (float64(sn)/sn_sum + 1), (float64(swn)/swn_sum + 1), a3, a4, a5)
	// 	}
	// }
}

func (info *Info) getAvg() float64 {
	if len(info.LeafList) == 0 {
		return 0
	}
	flen := float64(len(info.LeafList))
	sum := 0
	for _, l := range info.LeafList {
		sum += l
	}
	var sum2 float64 = 0
	avg := float64(sum) / flen
	for _, l := range info.LeafList {
		sum2 += (avg - float64(l)) * (avg - float64(l))
	}
	return math.Sqrt(sum2/flen + 1.0)
}
