package server

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/xunterr/indexp/indexer"
)

type Server struct {
	index *indexer.Index
}

type SearchResult struct {
	Filepath    string  `json:"filepath"`
	Score       float64 `json:"score"`
	Checksum    string  `json:"checksum"`
	LastIndexed string  `json:"indexedAt"`
}

func NewServer(index *indexer.Index) *Server {
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
			doc := s.index.Corpus[path]
			time := doc.IndexedAt
			results = append(results, SearchResult{
				Filepath:    path,
				Checksum:    doc.Checksum,
				LastIndexed: time.Format("2006-01-02 15:04:05"),
				Score:       score,
			})
		}
	}
	json.NewEncoder(w).Encode(results)
}

func (s *Server) Stats(w http.ResponseWriter, r *http.Request) {
	stats := struct {
		TotalDocs int
	}{len(s.index.Corpus)}
	t, err := template.ParseFiles("static/templates/index.html")
	if err != nil {
		log.Printf("Error executing template: %s", err.Error())
		return
	}
	err = t.Execute(w, stats)
	if err != nil {
		log.Printf("Error executing template: %s", err.Error())
	}
}
