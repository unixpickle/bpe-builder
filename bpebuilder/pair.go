package bpebuilder

import (
	"bytes"
	"hash/maphash"
)

var pairHash = maphash.MakeSeed()

type Pair struct {
	Left  []byte
	Right []byte
}

func (p Pair) Hash() uint64 {
	var h maphash.Hash
	h.SetSeed(pairHash)
	h.Write(p.Left)
	h.WriteByte(0)
	h.Write(p.Right)
	return h.Sum64()
}

func (p Pair) Equal(other Pair) bool {
	return bytes.Equal(p.Left, other.Left) && bytes.Equal(p.Right, other.Right)
}

type pairEntry[T any] struct {
	key   Pair
	value T
}

// A PairMap efficiently stores a map from Pair to any value.
type PairMap[T any] struct {
	mapping map[uint64][]pairEntry[T]
}

func NewPairMap[T any]() *PairMap[T] {
	return &PairMap[T]{mapping: map[uint64][]pairEntry[T]{}}
}

// AddToPairMap adds an integer value to an existing value, or sets the integer
// if the key was not present.
func AddToPairMap(p *PairMap[int], key Pair, addition int) {
	hash := key.Hash()
	if entry := p.getEntry(key, hash); entry != nil {
		entry.value += addition
	} else {
		p.mapping[hash] = append(p.mapping[hash], pairEntry[int]{key: key, value: addition})
	}
}

func (p *PairMap[T]) Set(key Pair, value T) {
	hash := key.Hash()
	newEntry := pairEntry[T]{key: key, value: value}
	if entry := p.getEntry(key, hash); entry != nil {
		*entry = newEntry
	} else {
		p.mapping[hash] = append(p.mapping[hash], newEntry)
	}
}

func (p *PairMap[T]) Get(key Pair) (T, bool) {
	if entry := p.getEntry(key, key.Hash()); entry != nil {
		return entry.value, true
	} else {
		var zero T
		return zero, false
	}
}

func (p *PairMap[T]) Delete(key Pair) (T, bool) {
	var zero T
	hash := key.Hash()
	records := p.mapping[hash]
	for i, x := range records {
		if x.key.Equal(key) {
			records[i] = records[len(records)-1]

			// Eliminate memory usage
			records[len(records)-1].value = zero
			records[len(records)-1].key = Pair{}

			p.mapping[hash] = records[:len(records)-1]
			return x.value, true
		}
	}
	return zero, false
}

func (p *PairMap[T]) Iterate(f func(key Pair, value T)) {
	for _, v := range p.mapping {
		for _, pairEntry := range v {
			f(pairEntry.key, pairEntry.value)
		}
	}
}

func (p *PairMap[T]) getEntry(key Pair, hash uint64) *pairEntry[T] {
	records := p.mapping[hash]
	for i, x := range records {
		if x.key.Equal(key) {
			return &records[i]
		}
	}
	return nil
}
