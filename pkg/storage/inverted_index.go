package storage

import (
	"math"
	"sort"
	"strings"
	"sync"
)

type Posting struct {
	DocID     uint64
	Frequency int
}

type InvertedIndex struct {
	mu           sync.RWMutex
	postings     map[string][]Posting
	docLengths   map[uint64]int
	totalDocs    int
	avgDocLength float64
}

func NewInvertedIndex() *InvertedIndex {
	return &InvertedIndex{
		postings:   make(map[string][]Posting),
		docLengths: make(map[uint64]int),
	}
}

func (idx *InvertedIndex) IndexDocument(docID uint64, text string) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	tokens := tokenize(text)
	idx.docLengths[docID] = len(tokens)
	idx.totalDocs++

	termFreq := make(map[string]int)
	for _, t := range tokens {
		termFreq[t]++
	}

	for term, freq := range termFreq {
		idx.postings[term] = append(idx.postings[term], Posting{
			DocID:     docID,
			Frequency: freq,
		})
	}
	idx.recalculateAvgDocLength()
}

func (idx *InvertedIndex) SearchBM25(queryStr string, topK int) []uint64 {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	queryTokens := tokenize(queryStr)
	scores := make(map[uint64]float64)

	k1 := 1.2
	b := 0.75

	for _, token := range queryTokens {
		postings, exists := idx.postings[token]
		if !exists {
			continue
		}

		df := len(postings)
		idf := math.Log(1.0 + (float64(idx.totalDocs)-float64(df)+0.5)/(float64(df)+0.5))

		for _, p := range postings {
			docLen := float64(idx.docLengths[p.DocID])
			tf := float64(p.Frequency)
			num := tf * (k1 + 1.0)
			den := tf + k1*(1.0-b+b*(docLen/idx.avgDocLength))
			score := idf * (num / den)
			scores[p.DocID] += score
		}
	}

	type docScore struct {
		id    uint64
		score float64
	}
	var scoredList []docScore
	for id, s := range scores {
		scoredList = append(scoredList, docScore{id: id, score: s})
	}
	sort.Slice(scoredList, func(i, j int) bool {
		return scoredList[i].score > scoredList[j].score
	})

	var results []uint64
	for i := 0; i < len(scoredList) && i < topK; i++ {
		results = append(results, scoredList[i].id)
	}
	return results
}

func (idx *InvertedIndex) recalculateAvgDocLength() {
	if idx.totalDocs == 0 {
		idx.avgDocLength = 0
		return
	}
	totalLen := 0
	for _, l := range idx.docLengths {
		totalLen += l
	}
	idx.avgDocLength = float64(totalLen) / float64(idx.totalDocs)
}

func tokenize(text string) []string {
	clean := strings.ToLower(text)
	words := strings.FieldsFunc(clean, func(r rune) bool {
		return (r < 'a' || r > 'z') && (r < '0' || r > '9')
	})
	return words
}
