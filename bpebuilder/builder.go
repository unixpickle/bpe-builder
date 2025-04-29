package bpebuilder

import (
	"bytes"
	"slices"

	"github.com/unixpickle/essentials"
)

const blockSize = 128

// CountPairs counts consecutive pairs of byte sequences.
func CountPairs(words [][][]byte) *PairMap[int] {
	allResult := NewPairMap[int]()
	essentials.ReduceConcurrentMap(0, len(words)/blockSize+1, func() (func(int), func()) {
		localResult := NewPairMap[int]()
		return func(blockIdx int) {
				for i := blockIdx * blockSize; i < (blockIdx+1)*blockSize && i < len(words); i++ {
					word := words[i]
					addWordCounts(localResult, word, 1)
				}
			}, func() {
				AddPairMap(allResult, localResult)
			}
	})
	return allResult
}

func addWordCounts(counts *PairMap[int], word [][]byte, multiplier int) {
	for j := 1; j < len(word); j++ {
		AddToPairMap(counts, Pair{Left: word[j-1], Right: word[j]}, multiplier)
	}
}

// Every time a pair is encountered in a word, replace it with the concatenated
// pair instead.
//
// Returns a map of the delta in pair counts before and after the merges.
func MergePairs(words [][][]byte, pair Pair) *PairMap[int] {
	combined := pair.Concat()

	outputDelta := NewPairMap[int]()
	essentials.ReduceConcurrentMap(0, len(words)/blockSize+1, func() (func(int), func()) {
		localDelta := NewPairMap[int]()
		return func(blockIdx int) {
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
					if newWord == nil {
						continue
					}
					addWordCounts(localDelta, word, -1)
					addWordCounts(localDelta, newWord, 1)
					words[i] = newWord
				}
			}, func() {
				AddPairMap(outputDelta, localDelta)
			}
	})
	return outputDelta
}
