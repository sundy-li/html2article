package html2article

import (
	"math"
)

type Info struct {
	TextCount     int
	LinkTextCount int
	TagCount      int
	LinkTagCount  int
	LeafList      []int
	Density       float64
	Pcount        int

	Data string

	score      float64
	DensitySum float64
}

func NewInfo() *Info {
	return &Info{}
}

func (info *Info) CalScore() {
	avg := info.getAvg()
	info.score = math.Log(avg) * float64(info.DensitySum) * math.Log(float64(info.TextCount-info.LinkTextCount+1)) * math.Log10(float64(info.Pcount+2))
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
