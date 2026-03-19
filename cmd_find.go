package main

import (
	"fmt"

	"github.com/amiv1/lsq-fork/system"
	"github.com/amiv1/lsq-fork/trie"
	"github.com/spf13/cobra"
)

var findOpenFlag bool

var findCmd = &cobra.Command{
	Use:          "find <prefix>",
	Short:        "Search page filenames by prefix",
	Aliases:      []string{"f"},
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return err
		}
		t, err := trie.Init(cfg.PagesDir)
		if err != nil {
			return fmt.Errorf("error loading pages directory: %w", err)
		}
		results := t.Search(args[0])
		if len(results) == 0 {
			fmt.Println("No results found")
			return nil
		}
		if findOpenFlag {
			system.LoadEditor(editorFlag, fmt.Sprintf("%s/%s", cfg.PagesDir, results[0]))
			return nil
		}
		fmt.Println("Search Results:")
		for _, r := range results {
			fmt.Println(r)
		}
		return nil
	},
}

func init() {
	findCmd.Flags().BoolVarP(&findOpenFlag, "open", "o", false, "Open the first matching page in editor")
}
