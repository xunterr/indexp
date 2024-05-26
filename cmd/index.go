/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"io/fs"
	"log"
	"path/filepath"

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
		err := Walk(index)
		if err != nil {
			log.Fatalf(err.Error())
			return
		}

		log.Println("Indexed! Saving...")
		err = utils.SaveJSON[indexer.Corpus]("corpus.json", index.GetCorpus())
		if err != nil {
			log.Fatalln(err.Error())
		}
		err = utils.SaveJSON[map[string]float64]("idf.json", index.GetIDFTable())
		if err != nil {
			log.Fatalln(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(indexCmd)

	// indexCmd.PersistentFlags().String("foo", "", "A help for foo")
	START_PATH = *indexCmd.Flags().StringP("path", "p", START_PATH, "starting point")
	IS_RECURSIVE = *indexCmd.Flags().BoolP("recursive", "r", IS_RECURSIVE, "starting point")
}

func Walk(index *indexer.Index) error {
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
