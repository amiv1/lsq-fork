package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jrswab/lsq/system"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:          "get",
	Short:        "Print journal or page to STDOUT",
	Aliases:      []string{"g"},
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runShowJournal("")
	},
}

var getYesterdayCmd = &cobra.Command{
	Use:          "yesterday",
	Short:        "Print yesterday's journal",
	Aliases:      []string{"y"},
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		date := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
		return runShowJournal(date)
	},
}

var getAgoCmd = &cobra.Command{
	Use:          "ago <n>",
	Short:        "Print journal from N days ago",
	Aliases:      []string{"a"},
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		n, err := strconv.Atoi(args[0])
		if err != nil || n < 0 {
			return fmt.Errorf("expected a non-negative integer for <n>, got %q", args[0])
		}
		date := time.Now().AddDate(0, 0, -n).Format("2006-01-02")
		return runShowJournal(date)
	},
}

var getDateCmd = &cobra.Command{
	Use:          "date <yyyy-MM-dd>",
	Short:        "Print journal for a specific date",
	Aliases:      []string{"d"},
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runShowJournal(args[0])
	},
}

var getPageCmd = &cobra.Command{
	Use:          "page <name>",
	Short:        "Print a specific page",
	Aliases:      []string{"p"},
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runShowPage(args[0])
	},
}

func init() {
	getCmd.AddCommand(getYesterdayCmd, getAgoCmd, getDateCmd, getPageCmd)
}

func runShowJournal(specDate string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := checkDirExists(cfg.DirPath); err != nil {
		return err
	}
	journalPath, err := system.GetJournal(cfg, cfg.JournalsDir, specDate)
	if err != nil {
		return fmt.Errorf("error resolving journal path: %w", err)
	}
	return system.PrintFile(journalPath)
}

func runShowPage(name string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	return system.PrintFile(resolvePagePath(cfg, name))
}
