/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"io/fs"
	"log"
	"path/filepath"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/xunterr/indexp/indexer"
	"github.com/xunterr/indexp/utils"
)

var (
	START_PATH   = "/home/xunterr/rfc/"
	IS_RECURSIVE = false
)

// indexCmd represents the index command
var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Index a directory",
	Run: func(cmd *cobra.Command, args []string) {
		index := indexer.NewEmptyIndex()
		var wg sync.WaitGroup
		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Prefix = "|---Indexing---|\t"
		s.Start()

		start := time.Now()
		err := Walk(&wg, START_PATH, index)
		wg.Wait()
		elapsed := time.Since(start)

		s.Stop()
		if err != nil {
			log.Fatalf(err.Error())
			return
		}

		log.Printf("Indexing ended in %d ms! Saving...", elapsed.Milliseconds())
		err = utils.SaveJSON[map[string]indexer.Document]("corpus.json", index.Corpus)
		if err != nil {
			log.Fatalln(err.Error())
		}
		err = utils.SaveJSON[map[string]float64]("idf.json", index.IdfTable)
		if err != nil {
			log.Fatalln(err.Error())
		}
	},
}

func init() {
	// indexCmd.PersistentFlags().String("foo", "", "A help for foo")
	indexCmd.Flags().StringVarP(&START_PATH, "path", "p", START_PATH, "starting point")
	indexCmd.Flags().BoolVarP(&IS_RECURSIVE, "recursive", "r", IS_RECURSIVE, "starting point")
	rootCmd.AddCommand(indexCmd)

}

func Walk(wg *sync.WaitGroup, start_path string, index *indexer.Index) error {
	return filepath.WalkDir(start_path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && IS_RECURSIVE {
			wg.Add(1)
			go func() {
				Walk(wg, path, index)
				wg.Done()
			}()
			return nil
		}

		index.IndexDoc(path)
		return nil
	})
}
