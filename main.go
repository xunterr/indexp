package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
)

var (
	START_PATH   = "./"
	IS_RECURSIVE = false
)

func init() {
	flag.StringVar(&START_PATH, "path", START_PATH, "starting point")
	flag.BoolVar(&IS_RECURSIVE, "recursive", IS_RECURSIVE, "recursive walking")
}

func main() {
	flag.Parse()
	index := NewIndex()
	err := Walk(index)
	if err != nil {
		log.Fatalf(err.Error())
		return
	}

	server := NewServer(index)
	http.HandleFunc("/search", server.Search)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln(err.Error())
		os.Exit(1)
	}

}

func Walk(index *Index) error {
	return filepath.WalkDir(START_PATH, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && IS_RECURSIVE {
			return Walk(index)
		}

		index.IndexDoc(path)
		return nil
	})
}

func DisplayFT(ft map[string]float64, num int) {
	keys := make([]string, 0, len(ft))

	for key := range ft {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return ft[keys[i]] > ft[keys[j]]
	})

	for i := 0; i < num && i < len(keys); i++ {
		key := keys[i]
		fmt.Printf("%d. %s -- %f \n", i+1, key, ft[key])
	}
}
