package utils

import (
	"github.com/spf13/cast"
	"math/rand"
)

// RateChoose 随机根据概率列表返回下标
func RateChoose(rateArr []float64) int {
	var selections []int
	var count = 0
	for _, rate := range rateArr {
		r := cast.ToInt(rate * 100)
		count += r
		selections = append(selections, count)
	}
	// 随机数
	randNum := rand.Intn(count)

	var index = 0
	for i, num := range selections {
		if randNum < num {
			index = i
			break
		}
	}
	return index
}

func RateSample(rateArr []float64, notBlank bool) []int {
	var samples []int
	for idx, rate := range rateArr {
		selected := RateChoose([]float64{100.0 - rate, rate})
		if selected == 1 { // 选中
			samples = append(samples, idx)
		}
	}
	// 如果必须
	if notBlank && len(samples) == 0 {
		return RateSample(rateArr, true)
	}
	return samples
}
