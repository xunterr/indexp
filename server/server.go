package server

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/xunterr/indexp/file"
	"github.com/xunterr/indexp/indexer"
)

type Server struct {
	index *indexer.Index
}

type SearchResult struct {
	Title       string  `json:"title"`
	Filepath    string  `json:"filepath"`
	Score       float64 `json:"score"`
	Checksum    string  `json:"checksum"`
	Snippet     string  `json:"snippet"`
	LastIndexed string  `json:"indexedAt"`
}

func NewServer(index *indexer.Index) *Server {
	return &Server{
		index: index,
	}
}

func (s *Server) Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	searchQuery := r.URL.Query().Get("query")
	maxQuery := r.URL.Query().Get("max")
	scores := s.index.Query(searchQuery)

	keys := make([]string, 0, len(scores))
	for key := range scores {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool { return scores[keys[i]].Score > scores[keys[j]].Score })

	var results []SearchResult

	if maxQuery != "" {
		maxResults, err := strconv.Atoi(maxQuery)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		if len(keys) > maxResults && maxResults >= 0 {
			keys = keys[:maxResults]
		}
	}
	for i, key := range keys {
		score := scores[key]
		if score.Score != 0 {
			doc := s.index.Corpus[key]
			time := doc.IndexedAt

			snippet := ""
			if i <= 5 {
				snippet, _ = getSnippet(key, score)
			}
			results = append(results, SearchResult{
				Title:       doc.Title,
				Filepath:    key,
				Checksum:    doc.Checksum,
				LastIndexed: time.Format("2006-01-02 15:04:05"),
				Score:       score.Score,
				Snippet:     snippet,
			})
		}
	}
	json.NewEncoder(w).Encode(results)
}

func getSnippet(path string, score indexer.Score) (string, error) {
	var lines []int
	for _, freq := range score.Tf {
		lines = append(lines, freq.FirstOccLine)
	}

	return file.GetSnippet(path, lines)
}

func (s *Server) IndexPage(w http.ResponseWriter, r *http.Request) {
	stats := struct {
		TotalDocs int
	}{len(s.index.Corpus)}
	t, err := template.ParseFiles("static/templates/index.html")
	if err != nil {
		log.Printf("Error parsing template: %s", err.Error())
		return
	}
	err = t.Execute(w, stats)
	if err != nil {
		log.Printf("Error executing template: %s", err.Error())
	}
}
