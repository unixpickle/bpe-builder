package bpebuilder

import (
	"bytes"
	"slices"

	"github.com/unixpickle/essentials"
)

const blockSize = 128

// CountPairs counts consecutive pairs of integers.
func CountPairs(words [][][]byte) *PairMap[int] {
	allResult := NewPairMap[int]()
	essentials.ReduceConcurrentMap(0, len(words)/blockSize+1, func() (func(int), func()) {
		localResult := NewPairMap[int]()
		return func(blockIdx int) {
				for i := blockIdx * blockSize; i < (blockIdx+1)*blockSize && i < len(words); i++ {
					word := words[i]
					for j := 1; j < len(word); j++ {
						AddToPairMap(localResult, Pair{Left: word[j-1], Right: word[j]}, 1)
					}
				}
			}, func() {
				localResult.Iterate(func(key Pair, value int) {
					AddToPairMap(allResult, key, value)
				})
			}
	})
	return allResult
}

func MergePairs(words [][][]byte, pair Pair) {
	combined := append(slices.Clone(pair.Left), pair.Right...)

	essentials.ConcurrentMap(0, len(words)/blockSize+1, func(blockIdx int) {
		for i := blockIdx * blockSize; i < (blockIdx+1)*blockSize && i < len(words); i++ {
			word := words[i]
			var newWord [][]byte
			for j := 0; j < len(word); j++ {
				if j+1 < len(word) &&
					bytes.Equal(word[j], pair.Left) && bytes.Equal(word[j+1], pair.Right) {
					if newWord == nil {
						newWord = slices.Clone(word[:j])
					}
					newWord = append(newWord, combined)
					j++
				} else if newWord != nil {
					newWord = append(newWord, word[j])
				}
			}
			if newWord != nil {
				words[i] = newWord
			}
		}
	})
}
