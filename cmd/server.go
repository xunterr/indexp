package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/xunterr/indexp/file"
	"github.com/xunterr/indexp/indexer"
	"github.com/xunterr/indexp/server"
)

var (
	PORT = 8080
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a server",
	Run: func(cmd *cobra.Command, args []string) {
		corpus := make(map[string]indexer.Document)
		err := file.LoadGOB[map[string]indexer.Document](INDEX_FILES[0], &corpus)
		if err != nil {
			log.Fatalln(err.Error())
		}

		docOcc := make(map[string]int)
		err = file.LoadGOB[map[string]int](INDEX_FILES[1], &docOcc)
		if err != nil {
			log.Fatalln(err.Error())
		}

		index := &indexer.Index{
			Corpus:         corpus,
			DocOccurrences: docOcc,
		}
		server := server.NewServer(index)
		http.HandleFunc("/search", server.Search)
		http.HandleFunc("/", server.IndexPage)
		http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("static/assets"))))

		err = http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
		if err != nil {
			log.Fatalln(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	PORT = *serverCmd.Flags().IntP("port", "p", PORT, "Port number")
}
