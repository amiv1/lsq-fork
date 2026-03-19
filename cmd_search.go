package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/jrswab/lsq/system"
	"github.com/spf13/cobra"
)

var searchOpenFlag bool

var searchCmd = &cobra.Command{
	Use:          "search <query>",
	Short:        "Search content for a keyword, or /pattern/ for regex",
	Aliases:      []string{"s"},
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		_, regexPattern := resolveSearch(args[0])
		pattern, err := regexp.Compile(regexPattern)
		if err != nil {
			return fmt.Errorf("error compiling search pattern: %w", err)
		}
		cfg, err := loadConfig()
		if err != nil {
			return err
		}
		dirs := []string{cfg.JournalsDir, cfg.PagesDir}
		if searchOpenFlag {
			for _, dir := range dirs {
				first, findErr := findFirstMatchInDirectory(dir, pattern)
				if findErr != nil {
					return findErr
				}
				if first != "" {
					system.LoadEditor(editorFlag, first)
					return nil
				}
			}
			fmt.Println("No results found")
			return nil
		}
		for _, dir := range dirs {
			if dirErr := searchInDirectory(dir, pattern); dirErr != nil {
				fmt.Fprintln(os.Stderr, "error searching directory:", dirErr)
			}
		}
		return nil
	},
}

func init() {
	searchCmd.Flags().BoolVarP(&searchOpenFlag, "open", "o", false, "Open the first matching file in editor")
}

// searchInFile prints lines matching pattern from filePath.
func searchInFile(filePath string, pattern *regexp.Regexp) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		if pattern.MatchString(line) {
			fmt.Printf("%s#%d: %s\n", filePath, lineNumber, line)
		}
	}
	return scanner.Err()
}

// searchInDirectory searches all files in directory for lines matching pattern.
func searchInDirectory(directory string, pattern *regexp.Regexp) error {
	return filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		return searchInFile(path, pattern)
	})
}

// findFirstMatchInDirectory returns the path of the first file containing a match.
func findFirstMatchInDirectory(directory string, pattern *regexp.Regexp) (string, error) {
	var firstFile string
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || firstFile != "" {
			return err
		}
		has, checkErr := fileHasMatch(path, pattern)
		if checkErr != nil {
			return checkErr
		}
		if has {
			firstFile = path
			return filepath.SkipAll
		}
		return nil
	})
	return firstFile, err
}

// fileHasMatch reports whether any line in filePath matches pattern.
func fileHasMatch(filePath string, pattern *regexp.Regexp) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if pattern.MatchString(scanner.Text()) {
			return true, nil
		}
	}
	return false, scanner.Err()
}
