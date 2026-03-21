package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/amiv1/lsq/config"
	"github.com/spf13/cobra"
	"golang.org/x/term"

	glamour "charm.land/glamour/v2"
	"charm.land/lipgloss/v2"
)

var noRawFlag bool

var getJournalCmd = &cobra.Command{
	Use:          "journal",
	Short:        "Browse all journal entries in a pager, newest first",
	Aliases:      []string{"j"},
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runOverview(noRawFlag)
	},
}

func init() {
	getJournalCmd.Flags().BoolVar(&noRawFlag, "raw", false, "Print raw Markdown without formatting or colours")
}

// runOverview streams formatted journal entries into a pager (less -R / $PAGER),
// falling back to stdout if no pager is available.
func runOverview(raw bool) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := checkDirExists(cfg.DirPath); err != nil {
		return err
	}

	files, err := journalFiles(cfg)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		fmt.Fprintln(os.Stderr, "No journal entries found.")
		return nil
	}

	// Attempt to start a pager. Fall back to stdout.
	out, cleanup, err := startPager()
	if err != nil {
		out = os.Stdout
		cleanup = func() {}
	}
	defer cleanup()

	goFmt := config.ConvertDateFormat(cfg.FileFmt)

	for _, name := range files {
		path := filepath.Join(cfg.JournalsDir, name)

		data, readErr := os.ReadFile(path)
		if readErr != nil || len(strings.TrimSpace(string(data))) == 0 {
			continue
		}

		// Parse date from filename (strip extension).
		stem := strings.TrimSuffix(name, filepath.Ext(name))
		t, parseErr := time.Parse(goFmt, stem)
		if parseErr != nil {
			t = time.Time{} // unknown date — still show it
		}

		header := dateHeader(t, raw)
		body := renderContent(string(data), raw)

		fmt.Fprintln(out, header)
		fmt.Fprintln(out, body)
	}

	return nil
}

// journalFiles returns journal filenames from the journals directory,
// sorted newest first (lexicographic descending, which works for date-named files).
func journalFiles(cfg *config.Config) ([]string, error) {
	entries, err := os.ReadDir(cfg.JournalsDir)
	if err != nil {
		return nil, fmt.Errorf("error reading journals directory: %w", err)
	}

	ext := ".md"
	if strings.EqualFold(cfg.FileType, "Org") {
		ext = ".org"
	}

	var names []string
	for _, e := range entries {
		if !e.IsDir() && strings.EqualFold(filepath.Ext(e.Name()), ext) {
			names = append(names, e.Name())
		}
	}

	sort.Sort(sort.Reverse(sort.StringSlice(names)))
	return names, nil
}

// dateHeader returns a styled (or plain) date header line.
func dateHeader(t time.Time, raw bool) string {
	var label string
	if t.IsZero() {
		label = "unknown date"
	} else {
		label = t.Format("2006-01-02 Mon")
	}

	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width <= 0 {
		width = 80
	}
	// Separator fills remaining width after "label " (label + space).
	sepLen := width - len(label) - 1
	if sepLen < 1 {
		sepLen = 1
	}
	separator := strings.Repeat("─", sepLen)

	if raw {
		return fmt.Sprintf("%s %s", label, separator)
	}

	dateStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)
	return fmt.Sprintf("%s %s", dateStyle.Render(label), separator)
}

var (
	linkRe = regexp.MustCompile(`\[\[([^\]]+)\]\]`)
	tagRe  = regexp.MustCompile(`(?:^|\s)(#\w+)`)

	// looseListPara matches a blank line (possibly whitespace-only) followed by
	// an indented continuation inside a list item. Glamour/goldmark fails to
	// insert a newline between the two paragraphs, causing text to run together.
	// Collapsing them into a tight list item fixes the rendering.
	looseListPara = regexp.MustCompile(`(?m)\n[ \t]*\n([ \t]+)`)

	linkStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
	tagStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
)

// normaliseForGlamour converts loose list item continuations into tight ones.
// This works around a glamour/goldmark bug where paragraphs inside a loose
// list item are rendered without any separator.
func normaliseForGlamour(content string) string {
	return looseListPara.ReplaceAllString(content, "\n$1")
}

// renderContent renders journal markdown content, applying glamour formatting
// and Logseq-specific colouring ([[links]], #tags). In raw mode it returns as-is.
func renderContent(content string, raw bool) string {
	if raw {
		return content
	}

	rendered, err := glamour.Render(normaliseForGlamour(content), "dark")
	if err != nil {
		// Fall back to unformatted content on render error.
		rendered = content
	}

	// Post-process glamour output to colour Logseq-specific tokens.
	rendered = linkRe.ReplaceAllStringFunc(rendered, func(m string) string {
		return linkStyle.Render(m)
	})
	rendered = tagRe.ReplaceAllStringFunc(rendered, func(m string) string {
		// Preserve any leading whitespace before the #tag.
		idx := strings.Index(m, "#")
		return m[:idx] + tagStyle.Render(m[idx:])
	})

	return rendered
}

// startPager starts $PAGER (defaulting to "less -R") and returns a writer
// connected to its stdin, plus a cleanup function that closes the pipe and
// waits for the pager to exit.
func startPager() (io.Writer, func(), error) {
	pagerCmd := os.Getenv("PAGER")
	if pagerCmd == "" {
		pagerCmd = "less"
	}

	// Split into command + args (e.g. "less -R").
	parts := strings.Fields(pagerCmd)
	if pagerCmd == "less" {
		parts = []string{"less", "-R"}
	}

	cmd := exec.Command(parts[0], parts[1:]...) //nolint:gosec
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	pipe, err := cmd.StdinPipe()
	if err != nil {
		return nil, nil, err
	}

	if err := cmd.Start(); err != nil {
		pipe.Close()
		return nil, nil, err
	}

	cleanup := func() {
		pipe.Close()
		cmd.Wait() //nolint:errcheck
	}

	return pipe, cleanup, nil
}
