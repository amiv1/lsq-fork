<!--<p align="center">
<img width="25%" src="https://lsq.sh/media/img/lsq_logo_cropped.png" alt="lsq logo">
</p>-->

# lsq

[![Go Report Card](https://goreportcard.com/badge/github.com/amiv1/lsq-fork)](https://goreportcard.com/report/github.com/amiv1/lsq-fork)
[![Release](https://img.shields.io/github/v/release/amiv1/lsq-fork)](https://github.com/amiv1/lsq-fork/releases)
[![License](https://img.shields.io/github/license/amiv1/lsq-fork)](https://github.com/amiv1/lsq-fork/blob/master/LICENSE)

The ultra-fast CLI companion for [Logseq](https://github.com/logseq/logseq) designed to speed up your note capture directly from the terminal!

> **Fork of [jrswab/lsq](https://github.com/jrswab/lsq)** — aimed at improving UX and adding more features.

## Why lsq?
- ⚡️ Lightning-fast journal additions without leaving your terminal
- ⌨️ Optimized for both quick captures and extended writing sessions
- 🎯 Native support for Logseq's file naming and formatting conventions
- 🔄 Seamless integration with your existing Logseq workflow
- 💻 Built by Logseq users, for Logseq users

## Features That Power Your Workflow
- External editor integration ($EDITOR by default)
- Automatic journal file creation
- Support for both Markdown and Org formats
- Configurable file naming format
- Customizable directory location with `~` and environment variable support
- STDOUT output mode for shell scripting and widget integration
- User defined configuration file

## Installation

If you have the [original lsq](https://github.com/jrswab/lsq) installed with `go install`, remove it first:

```shell
rm $(which lsq)
```

### macOS — Homebrew

```shell
brew install amiv1/lsq/lsq
```

### Debian / Ubuntu

Add a repository:
```shell
curl -fsSL https://apt.fury.io/amiv1/gpg.key | sudo gpg --dearmor -o /etc/apt/trusted.gpg.d/fury-amiv1.gpg
echo "deb [signed-by=/etc/apt/trusted.gpg.d/fury-amiv1.gpg] https://apt.fury.io/amiv1/ * *" | sudo tee /etc/apt/sources.list.d/lsq-fork.list
```

Install lsq:
```shell
sudo apt update && sudo apt install lsq
```

### Fedora / RHEL / CentOS

Add a repository:
```shell
echo "[fury-amiv1]
name=lsq (Gemfury)
baseurl=https://yum.fury.io/amiv1/
enabled=1
gpgcheck=0" | sudo tee /etc/yum.repos.d/fury-amiv1.repo
```

Install lsq:
```shell
sudo dnf install lsq
```

### Build and install from source

Requires Go:

```bash
git clone https://github.com/amiv1/lsq-fork.git
cd lsq-fork
go install .
```

Make sure you have the location of the Go binaries in your `$PATH`. Run go env and find the variable called GOPATH.
Then copy that location to your shell's `$PATH` if it's not already there.

## Usage

### Commands

Each command has a short alias. All commands accept `--dir` and `--editor` as global flags.

#### Journal commands
Open the journal in your editor, or append text if provided.

| Command | Alias | Description |
|---------|-------|-------------|
| `lsq today [text...]` | `lsq t` | Today's journal |
| `lsq yesterday [text...]` | `lsq y` | Yesterday's journal |
| `lsq ago <n> [text...]` | `lsq a <n>` | Journal from N days ago |

When piped (`echo "text" | lsq t`), STDIN is appended automatically.  
Flag: `-i`/`--indent <n>` — indentation level for appended text.

#### `lsq page <name> [text...]` / `lsq p`
Open a specific page in your editor, or append text to it. File extension is auto-detected.

Flag: `-i`/`--indent <n>`

#### `lsq get` / `lsq g`
Print to STDOUT. Uses subcommands for date/target selection (default: today).

| Subcommand | Alias | Example |
|------------|-------|---------|
| *(default)* | — | `lsq g` |
| `get yesterday` | `g y` | `lsq g y` |
| `get ago <n>` | `g a <n>` | `lsq g a 3` |
| `get date <yyyy-MM-dd>` | `g d <date>` | `lsq g d 2024-01-15` |
| `get page <name>` | `g p <name>` | `lsq g p notes` |

#### `lsq search <query>` / `lsq s`
Search file **contents** across all journals and pages. Plain text matches literally; wrap in `/…/` for regex.

Flag: `-o`/`--open` — open the first matching file in editor.

#### `lsq find <prefix>` / `lsq f`
Search page **filenames** by prefix.

Flag: `-o`/`--open` — open the first matching page in editor.

#### Global flags
Available on all commands:
- `--dir <path>` — Logseq directory (overrides config). Supports `~` and environment variables.
- `--editor <cmd>` — Editor to use (default: `$EDITOR`, fallback: vim).

### Configuration File
This file must be stored in your config directory as `lsq/config.edn`.
On Unix systems, it returns `$XDG_CONFIG_HOME` if non-empty, else `$HOME/.config` will be used.
On macOS, it returns `$HOME/Library/Application Support`.
On Windows, it returns `%AppData%`.
On Plan 9, it returns `$home/lib`.

#### Configuration Behavior
The configuration file will override any lsq defaults which are defined. If a CLI flag is provided, the flag value will override the config file value.

#### Configuration File Example:
```EDN
{
  ;; Either "Markdown" or "Org".
  :file/type "Markdown"
  ;; This will be used for journal file names
  ;; Using the format below and the file type above will produce 2025.01.01.md
  :file/format "yyyy_MM_dd"
  ;; The directory which holds all your notes
  ;; Supports ~ and environment variables (e.g., ~/Logseq or $HOME/Logseq)
  :directory "~/Logseq"
}
```
**Note:** The configured directory must contain both a `journals` and `pages` subdirectory for lsq to function properly. These are automatically created when using Logseq, but will need to be manually created if setting lsq to use a new directory or without Logseq.

### Usage Examples:

```shell
lsq t
```
Opens today's journal in `$EDITOR`.

```shell
lsq t "Entry text here"
```
Appends `Entry text here` as a bullet point to today's journal.

```shell
lsq a 2 "Entry text here"
```
Appends to the journal from 2 days ago.

```shell
lsq y "Entry text here"
```
Appends to yesterday's journal.

```shell
lsq p my-page "Entry text here"
```
Appends to the page named `my-page` (extension auto-detected). Without text, opens the page in editor.

```shell
lsq t --indent 1 "sub-item text"
```
Appends text as an indented bullet (one tab level deep). Use `--indent 2` for two levels, and so on.

```shell
lsq search TODO
```
Searches all journals and pages for lines containing `TODO`.

```shell
lsq search '/TODO|FIXME/'
```
Searches using the regex `TODO|FIXME`.

```shell
lsq search TODO --open
```
Searches for `TODO` and opens the first matching file in editor.

```shell
lsq find go
```
Lists pages whose filename starts with `go`.

```shell
lsq g
```
Prints today's journal to STDOUT. Useful for shell integration, piping, or display widgets.

```shell
lsq g a 3
```
Prints the journal from 3 days ago to STDOUT.

```shell
lsq g p notes
```
Prints the `notes` page to STDOUT.

```shell
cat ~/.zshrc | lsq t
```
Appends the contents of `~/.zshrc` to today's journal via STDIN.

```shell
run_long_batch_job |& lsq p "long-job.$(date +%s).log"
```
Appends STDOUT and STDERR of a long-running job to a new page.

## Contributing
For information on contributing to lsq check out [CONTRIBUTING.md](https://github.com/amiv1/lsq-fork/blob/master/CONTRIBUTING.md).

## License
GPL v3
