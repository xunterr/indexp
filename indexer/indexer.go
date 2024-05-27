package indexer

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"math"
	"path/filepath"
	"time"

	"github.com/xunterr/indexp/tokenizer"
)

type Document struct {
	Checksum  string
	Tf        map[string]float64
	IndexedAt time.Time
}

type Index struct {
	Corpus        map[string]Document
	DocOccurences map[string]int
}

func NewEmptyIndex() *Index {
	return &Index{
		Corpus:        make(map[string]Document),
		DocOccurences: make(map[string]int),
	}
}

func (index *Index) IndexDoc(path string, data []byte) Document {
	checksum := md5.Sum(data)
	tokenizer := tokenizer.NewTokenizer(data)
	occurences := make(map[string]int)
	for {
		token, err := tokenizer.ScanToken()
		if err != nil {
			if err == io.EOF {
				break
			}
			return Document{}
		}
		if token != "" {
			if _, ok := occurences[token]; !ok {
				docOcc := index.DocOccurences[token]
				index.DocOccurences[token] = docOcc + 1
			}

			freq := occurences[token]
			occurences[token] = freq + 1
		}
	}
	tf := CalcDocTF(occurences)
	doc := Document{
		Checksum:  hex.EncodeToString(checksum[:]),
		Tf:        tf,
		IndexedAt: time.Now(),
	}

	abs, _ := filepath.Abs(path)
	index.Corpus[abs] = doc
	return doc
}

func (index Index) Query(request string) map[string]float64 {
	tokenizer := tokenizer.NewTokenizer([]byte(request))
	tokens, err := tokenizer.ScanAll()
	if err != nil {
		return nil
	}

	score := make(map[string]float64)
	for path, doc := range index.Corpus {
		docScore := float64(0)
		for _, token := range tokens {
			if tf, ok := doc.Tf[token]; ok {
				docScore += tf * index.termIDF(token)
			}
		}
		score[path] = docScore
	}
	return score
}

func CalcDocTF(occurences map[string]int) map[string]float64 {
	tf := make(map[string]float64)
	tokensNum := 0
	for _, occNum := range occurences {
		tokensNum += occNum
	}

	for token, occNum := range occurences {
		if _, ok := tf[token]; !ok {
			tf[token] = float64(occNum) / float64(tokensNum)
		}
	}
	return tf
}

func (i Index) termIDF(term string) float64 {
	docOcc := i.DocOccurences[term]
	return math.Log(float64(len(i.Corpus)) / float64(docOcc))
}

func (i Index) getUniqueTerms() []string {
	var unique []string
	for _, doc := range i.Corpus {
		for term := range doc.Tf {
			unique = append(unique, term)
		}
	}
	return unique
}
