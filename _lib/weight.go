package _lib

import (
	"math/rand"
)

type WeightItem struct {
	Item   interface{} // 项目
	Weight int         // 权重
}

// WeightSelector 权重选择器，返回：WeightItem.Item
func WeightSelector(items []WeightItem) interface{} {
	totalWeight := 0
	for _, item := range items {
		totalWeight += item.Weight
	}
	r := rand.Int() * totalWeight
	var cumulativeWeight int
	for _, item := range items {
		cumulativeWeight += item.Weight
		if r <= cumulativeWeight {
			return item.Item
		}
	}
	return items[0].Item
}
