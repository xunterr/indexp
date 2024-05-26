package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	index *Index
}

type SearchResult struct {
	Filepath    string  `json:"filepath"`
	Score       float64 `json:"score"`
	Checksum    string  `json:"checksum"`
	LastIndexed string  `json:"indexedAt"`
}

func NewServer(index *Index) *Server {
	return &Server{
		index: index,
	}
}

func (s *Server) Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query().Get("query")
	scores := s.index.Query(query)
	var results []SearchResult
	for path, score := range scores {
		if score != 0 {
			doc := s.index.GetDoc(path)
			time := doc.IndexedAt
			results = append(results, SearchResult{
				Filepath: path,
				Checksum: doc.Checksum,
				LastIndexed: fmt.Sprintf("%d-%d-%d %d:%d",
					time.Day(), time.Month(), time.Year(), time.Hour(), time.Minute()),
				Score: score,
			})
		}
	}
	log.Printf("%v", results)
	json.NewEncoder(w).Encode(results)
}
