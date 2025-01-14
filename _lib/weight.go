package _lib

import (
	"math/rand"
	"sort"
)

type WeightItem struct {
	Item   interface{} // 项目
	Weight int         // 权重
}

// WeightSelector 权重选择器，返回：WeightItem.Item
func WeightSelector(items []WeightItem) interface{} {
	itemArray := make([]int, len(items))
	count := 0
	for i, item := range items {
		count += item.Weight
		itemArray[i] = count
	}
	r := rand.Intn(count) + 1
	pos := sort.SearchInts(itemArray, r)
	if pos < len(itemArray) {
		return items[pos].Item
	}
	return items[0].Item
}
