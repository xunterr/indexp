package indexer

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"math"
	"path/filepath"
	"time"

	"github.com/xunterr/indexp/file"
	"github.com/xunterr/indexp/tokenizer"
)

type Document struct {
	Title     string
	Checksum  string
	Tf        map[string]TermFreq
	IndexedAt time.Time
}

type TermFreq struct {
	Freq         float64
	FirstOccLine int
}

type Score struct {
	Score float64
	Tf    map[string]TermFreq
}

type TermOccurrences struct {
	Num          int
	FirstOccLine int
}

type Index struct {
	Corpus         map[string]Document
	DocOccurrences map[string]int
}

func NewEmptyIndex() *Index {
	return &Index{
		Corpus:         make(map[string]Document),
		DocOccurrences: make(map[string]int),
	}
}

func (index *Index) IndexDoc(path string, file *file.File) {
	checksum := md5.Sum(file.Data)
	tokensNum, occurrences := index.toTokens(file.Data)

	tf := make(map[string]TermFreq, len(occurrences))
	for token, occ := range occurrences {
		docOcc := index.DocOccurrences[token]
		index.DocOccurrences[token] = docOcc + 1

		tf[token] = TermFreq{
			Freq:         calcTF(occ.Num, tokensNum),
			FirstOccLine: occ.FirstOccLine,
		}
	}

	abs, _ := filepath.Abs(path)
	index.Corpus[abs] = Document{
		Title:     file.Title,
		Checksum:  hex.EncodeToString(checksum[:]),
		Tf:        tf,
		IndexedAt: time.Now(),
	}
}

func (index Index) toTokens(data []byte) (tokensNum int, occurrences map[string]TermOccurrences) {
	occurrences = map[string]TermOccurrences{}
	tokenizer := tokenizer.NewTokenizer(data)

	for {
		token, err := tokenizer.ScanToken()
		if err != nil {
			if err == io.EOF {
				break
			}
			continue
		}
		if token.Literal != "" {
			if occ, ok := occurrences[token.Literal]; !ok {
				occurrences[token.Literal] = TermOccurrences{1, token.Line}
			} else {
				occ.Num++
				occurrences[token.Literal] = occ
			}
			tokensNum++
		}
	}

	return
}

func (index Index) Query(request string) map[string]Score {
	tokenizer := tokenizer.NewTokenizer([]byte(request))
	tokens, err := tokenizer.ScanAll()
	if err != nil {
		return nil
	}

	scores := make(map[string]Score)
	for path, doc := range index.Corpus {
		var docScore float64
		reqTf := make(map[string]TermFreq)
		for _, token := range tokens {
			if tf, ok := doc.Tf[token.Literal]; ok {
				docScore += tf.Freq * index.calcDocIdf(token.Literal)
				reqTf[token.Literal] = tf
			}
		}
		scores[path] = Score{
			Score: docScore,
			Tf:    reqTf,
		}
	}
	return scores
}

func calcTF(termOcc int, termsNum int) float64 {
	return float64(termOcc) / float64(termsNum)
}

func (i Index) calcDocIdf(term string) float64 {
	docOcc := i.DocOccurrences[term]
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
