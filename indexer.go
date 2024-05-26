package main

import (
	"crypto/md5"
	"errors"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"time"
)

type Document struct {
	Checksum  string
	Tf        map[string]float64
	IndexedAt time.Time
}

type Corpus map[string]Document

type Index struct {
	corpus   Corpus
	idfTable map[string]float64
}

func NewIndex() *Index {
	return &Index{
		corpus:   make(Corpus),
		idfTable: make(map[string]float64),
	}
}

func (index *Index) IndexDoc(path string) Document {
	data, err := ReadFile(path)
	if err != nil {
		return Document{}
	}
	checksum := md5.Sum(data)
	tokenizer := NewTokenizer(data)
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
			freq := occurences[token]
			occurences[token] = freq + 1
		}
	}
	tf := CalcDocTF(occurences)
	doc := Document{
		Checksum:  string(checksum[:]),
		Tf:        tf,
		IndexedAt: time.Now(),
	}

	abs, _ := filepath.Abs(path)
	index.corpus[abs] = doc
	index.idfTable = index.corpus.CalculateIDF()
	return doc
}

func (index Index) GetDoc(path string) Document {
	return index.corpus[path]
}

func (index Index) Query(request string) map[string]float64 {
	tokenizer := NewTokenizer([]byte(request))
	tokens, err := tokenizer.ScanAll()
	if err != nil {
		return nil
	}

	score := make(map[string]float64)
	for path, doc := range index.corpus {
		docScore := float64(0)
		for _, token := range tokens {
			if tf, ok := doc.Tf[token]; ok {
				idf, ok := index.idfTable[token]
				if !ok {
					continue
				}
				docScore += tf * idf
			}
		}
		score[path] = docScore
	}
	return score
}

func (c Corpus) CalculateIDF() map[string]float64 {
	unique := c.GetUniqueTerms()
	docFreq := make(map[string]float64)
	for _, term := range unique {
		docFreq[term] = c.GetTermIDF(term)
	}
	return docFreq
}

func ReadFile(filename string) ([]byte, error) {
	fd, err := os.Open(filename)
	defer fd.Close()
	if err != nil {
		return nil, err
	}

	if hasExt(fd.Name(), []string{".txt", ".csv", ".md", ".json"}) {
		data, err := io.ReadAll(fd)
		if err != nil {
			log.Fatalf("Error reading file: %s", err.Error())
			return nil, err
		}
		return data, nil
	} else {
		return nil, errors.New("File extension is not supported! Skipping")
	}
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

func hasExt(path string, extensions []string) (ok bool) {
	hasExt := false
	fileExt := filepath.Ext(path)
	for _, ext := range extensions {
		if hasExt {
			break
		}
		hasExt = ext == fileExt
	}

	return hasExt
}

func TermFreqInSet(term string, occurences map[string]uint) float64 {
	var termsNum uint = 0
	for _, v := range occurences {
		termsNum += v
	}

	freq, _ := occurences[term]
	normalizedFreq := float64(freq) / float64(termsNum)
	return normalizedFreq
}

func (c Corpus) GetTermIDF(term string) float64 {
	docFreq := 0
	for _, doc := range c {
		if _, ok := doc.Tf[term]; ok {
			docFreq++
		}
	}

	return math.Log(float64(len(c)) / float64(docFreq))
}

func (c Corpus) GetUniqueTerms() []string {
	var unique []string
	for _, doc := range c {
		for term := range doc.Tf {
			unique = append(unique, term)
		}
	}
	return unique
}
