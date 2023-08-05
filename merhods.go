package main

import (
	"math/rand"
	"time"
)

func GetRandomWakas(slice []Waka, numElements int) []Waka {
	if numElements <= 0 || numElements > len(slice) {
		// 要素数が0以下または元のスライスよりも大きい場合は空のスライスを返す
		return []Waka{}
	}

	source := rand.NewSource(time.Now().UnixNano())
	rand.New(source)

	result := make([]Waka, numElements)

	// ランダムなインデックスを生成して、対応する要素を新しいスライスに追加
	for i := 0; i < numElements; i++ {
		randomIndex := rand.Intn(len(slice))
		result[i] = slice[randomIndex]

		// 重複しないように選択済みの要素をスライスから削除
		// 削除することでスライスが詰められるので、ランダムな選択が保証されます
		slice = append(slice[:randomIndex], slice[randomIndex+1:]...)
	}

	return result
}