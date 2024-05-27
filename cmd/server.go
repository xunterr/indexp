package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/xunterr/indexp/indexer"
	"github.com/xunterr/indexp/server"
	"github.com/xunterr/indexp/utils"
)

var (
	PORT = 8080
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a server",
	Run: func(cmd *cobra.Command, args []string) {
		corpus, err := utils.LoadJSON[map[string]indexer.Document]("corpus.json")
		if err != nil {
			log.Fatalln(err.Error())
		}

		idf, err := utils.LoadJSON[map[string]float64]("idf.json")
		if err != nil {
			log.Fatalln(err.Error())
		}

		index := &indexer.Index{
			Corpus:   *corpus,
			IdfTable: *idf,
		}
		server := server.NewServer(index)
		http.HandleFunc("/search", server.Search)
		http.Handle("/", http.FileServer(http.Dir("./static")))

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
