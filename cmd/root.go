/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
)

var (
	INDEX_FILES []string = []string{"corpus.json", "do.json"}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "indexp",
	Short: "Indexp is an indexing file search",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if len(INDEX_FILES) != 2 {
			return errors.New("Length of index files should be 2")
		}
		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringSliceVarP(&INDEX_FILES, "files", "f", INDEX_FILES, "index files (2)")
	rootCmd.AddCommand(indexCmd)

}
