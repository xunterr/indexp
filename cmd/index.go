package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/signal"
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
	Run:   Run,
}

func Run(cmd *cobra.Command, args []string) {
	index := indexer.NewEmptyIndex()
	var wg sync.WaitGroup

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Prefix = fmt.Sprintf("| %s |\t", START_PATH)

	end := make(chan struct{})
	go func() {
		t := time.NewTicker(1 * time.Second)
		for {
			s.Suffix = fmt.Sprintf("\t| Docs indexed: %d |", len(index.Corpus))
			select {
			case <-end:
				t.Stop()
				return
			case <-t.C:
				continue
			}
		}
	}()
	s.Start()

	start := time.Now()
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	err := Walk(&wg, ctx, START_PATH, index)
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
	err = utils.SaveJSON[map[string]int]("do.json", index.DocOccurences)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func init() {
	// indexCmd.PersistentFlags().String("foo", "", "A help for foo")
	indexCmd.Flags().StringVarP(&START_PATH, "path", "p", START_PATH, "starting point")
	indexCmd.Flags().BoolVarP(&IS_RECURSIVE, "recursive", "r", IS_RECURSIVE, "starting point")
	rootCmd.AddCommand(indexCmd)

}

func Walk(wg *sync.WaitGroup, ctx context.Context, start_path string, index *indexer.Index) error {
	return filepath.WalkDir(start_path, func(path string, d fs.DirEntry, err error) error {
		select {
		case <-ctx.Done():
			return filepath.SkipAll
		default:
		}
		if err != nil {
			return err
		}
		if d.IsDir() && IS_RECURSIVE {
			wg.Add(1)
			go func() {
				Walk(wg, ctx, path, index)
				wg.Done()
			}()
			return nil
		}

		data, err := ReadFile(path)
		if err == nil {
			index.IndexDoc(path, data)
		}

		return nil
	})
}

func ReadFile(filename string) ([]byte, error) {
	fd, err := os.Open(filename)
	defer fd.Close()
	if err != nil {
		return nil, err
	}

	buff := bufio.NewReader(fd)
	if hasExt(fd.Name(), []string{".txt", ".csv", ".md", ".json", ".xml"}) {
		data, err := io.ReadAll(buff)
		if err != nil {
			log.Fatalf("Error reading file: %s", err.Error())
			return nil, err
		}
		return data, nil
	} else {
		return nil, errors.New("File extension is not supported! Skipping")
	}
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
