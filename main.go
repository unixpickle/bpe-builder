package main

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"os"
	"sync"

	"github.com/dlclark/regexp2"
	"github.com/unixpickle/bpe-builder/bpebuilder"
	"github.com/unixpickle/essentials"
)

const DefaultRegexPattern = `'(?i:[sdmt]|ll|ve|re)|[^\r\n\p{L}\p{N}]?\p{L}+|\p{N}{1,3}| ?[^\s\p{L}\p{N}]+[\r\n]*|\s+$|\s*[\r\n]|\s+(?!\S)|\s`

func main() {
	var dataDir string
	var outputPath string
	var vocabSize int
	var pattern string
	flag.StringVar(&dataDir, "data-dir", "", "path to directory of JSON files")
	flag.StringVar(&outputPath, "output-path", "", "path to output vocabulary")
	flag.StringVar(&pattern, "pattern", DefaultRegexPattern, "regexp pattern to split up words")
	flag.IntVar(&vocabSize, "vocab-size", 8192, "the size of the vocabulary to produce")
	flag.Parse()

	if dataDir == "" || outputPath == "" {
		essentials.Die("You must specify -data-dir and -output-path. See -help.")
	}

	log.Println("loading data ...")
	data, err := bpebuilder.LoadJSONData(dataDir)
	essentials.Must(err)

	log.Println("splitting up words ...")
	re, err := regexp2.Compile(pattern, regexp2.None)
	essentials.Must(err)

	var allWords [][][]byte
	var lock sync.Mutex
	essentials.ConcurrentMap(0, len(data), func(i int) {
		var result [][]byte
		matches, err := FindAllString(re, data[i])
		essentials.Must(err)
		for _, match := range matches {
			for _, x := range []byte(match) {
				result = append(result, []byte{x})
			}
		}
		lock.Lock()
		allWords = append(allWords, result)
		lock.Unlock()
	})

	log.Println("building merges ...")
	counts := bpebuilder.CountPairs(allWords)
	vocab := make([][]byte, 256, vocabSize)
	for i := range 256 {
		vocab[i] = []byte{byte(i)}
	}
	for len(vocab) < vocabSize {
		maxCount := -1
		maxPair := bpebuilder.Pair{}
		counts.Iterate(func(key bpebuilder.Pair, value int) {
			if value > maxCount {
				maxCount = value
				maxPair = key
			}
		})
		log.Printf(" - at vocab size %d, max pair is %v with count %d", len(vocab), string(maxPair.Concat()), maxCount)
		delta := bpebuilder.MergePairs(allWords, maxPair)
		bpebuilder.AddPairMap(counts, delta)
		vocab = append(vocab, maxPair.Concat())
	}

	log.Printf("saving vocabulary to %s ...", outputPath)
	outData, err := json.Marshal(vocab)
	essentials.Must(err)
	essentials.Must(os.WriteFile(outputPath, outData, 0644))
}

func FindAllString(re *regexp2.Regexp, text string) ([]string, error) {
	var matches []string
	start := 0

	runes := []rune(text)

	for {
		match, err := re.FindRunesMatchStartingAt(runes, start)
		if err != nil {
			return nil, err
		}
		if match == nil {
			break
		}

		matches = append(matches, match.String())
		start = match.Index + match.Length

		if match.Length == 0 {
			return nil, errors.New("zero-length match")
		}
	}

	return matches, nil
}
