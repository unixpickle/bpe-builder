package bpebuilder

import (
	"bytes"
	"testing"
)

func TestCountPairs(t *testing.T) {
	allWords := [][][]byte{}
	for range 10000 {
		allWords = append(allWords, [][]byte{[]byte("h"), []byte("ij"), []byte("k")})
	}
	for range 1000 {
		allWords = append(allWords, [][]byte{[]byte("a"), []byte("ij"), []byte("k")})
	}
	counts := CountPairs(allWords)
	if n, ok := counts.Get(Pair{Left: []byte("a"), Right: []byte("ij")}); !ok || n != 1000 {
		t.Fatalf("got %d", n)
	}
	if n, ok := counts.Get(Pair{Left: []byte("h"), Right: []byte("ij")}); !ok || n != 10000 {
		t.Fatalf("got %d", n)
	}
	if n, ok := counts.Get(Pair{Left: []byte("ij"), Right: []byte("k")}); !ok || n != 11000 {
		t.Fatalf("got %d", n)
	}
}

func TestMergePairs(t *testing.T) {
	words := [][][]byte{
		{
			[]byte("he"),
			[]byte("llo"),
		},
		{
			[]byte("hi"),
			[]byte(" there"),
		},
		{
			[]byte("he"),
			[]byte("llo"),
			[]byte("llo"),
			[]byte("he"),
			[]byte("llo"),
			[]byte("he"),
			[]byte("he"),
			[]byte("llo"),
		},
		{
			[]byte("hi"),
			[]byte("he"),
			[]byte("llo"),
			[]byte("x"),
		},
	}
	delta := MergePairs(words, Pair{Left: []byte("he"), Right: []byte("llo")})

	if n, ok := delta.Get(Pair{Left: []byte("he"), Right: []byte("llo")}); !ok || n != -5 {
		t.Errorf("bad delta: %d", n)
	}
	if n, ok := delta.Get(Pair{Left: []byte("hello"), Right: []byte("llo")}); !ok || n != 1 {
		t.Errorf("bad delta: %d", n)
	}
	if n, ok := delta.Get(Pair{Left: []byte("hi"), Right: []byte("hello")}); !ok || n != 1 {
		t.Errorf("bad delta: %d", n)
	}
	if n, ok := delta.Get(Pair{Left: []byte("hello"), Right: []byte("x")}); !ok || n != 1 {
		t.Errorf("bad delta: %d", n)
	}
	if n, ok := delta.Get(Pair{Left: []byte("llo"), Right: []byte("llo")}); !ok || n != -1 {
		t.Errorf("bad delta: %d", n)
	}

	if len(words[0]) != 1 {
		t.Fatalf("%v", words[0])
	}
	if !bytes.Equal(words[0][0], []byte("hello")) {
		t.Fatalf("%v", words[0])
	}

	if len(words[1]) != 2 {
		t.Fatalf("%v", words[0])
	}
	if !bytes.Equal(words[1][0], []byte("hi")) || !bytes.Equal(words[1][1], []byte(" there")) {
		t.Fatalf("%v", words[1])
	}

	if len(words[2]) != 5 {
		t.Fatalf("%v", words[2])
	}
	for i, x := range []string{"hello", "llo", "hello", "he", "hello"} {
		if !bytes.Equal(words[2][i], []byte(x)) {
			var printMe []string
			for _, a := range words[2] {
				printMe = append(printMe, string(a))
			}
			t.Fatalf("%v", printMe)
		}
	}

	if len(words[3]) != 3 {
		t.Fatalf("%v", words[3])
	}
	for i, x := range []string{"hi", "hello", "x"} {
		if !bytes.Equal(words[3][i], []byte(x)) {
			var printMe []string
			for _, a := range words[3] {
				printMe = append(printMe, string(a))
			}
			t.Fatalf("%v", printMe)
		}
	}
}
