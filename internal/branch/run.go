package branch

import (
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/lhlyu/gitx/internal/git"
	"github.com/lhlyu/gitx/internal/repo"
	"github.com/lhlyu/gitx/internal/term"
)

var (
	titleColor   = color.New(color.FgCyan, color.Bold)
	currentColor = color.New(color.FgGreen, color.Bold)
	localColor   = color.New(color.FgYellow)
	remoteColor  = color.New(color.FgMagenta)
	hashColor    = color.New(color.FgYellow)
	dateColor    = color.New(color.FgWhite)
	infoColor    = color.New(color.FgWhite)
	errorColor   = color.New(color.FgRed, color.Bold)
)

type Entry struct {
	Name      string
	Kind      string
	Hash      string
	Date      string
	Subject   string
	IsCurrent bool
}

func Run() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	if !repo.IsGitRepo(currentDir) {
		_, _ = errorColor.Println("❌ 当前目录不是 Git 项目")
		return nil
	}

	client := git.NewClient()
	out, err := client.RunInDir(
		currentDir,
		"for-each-ref",
		"--format=%(HEAD)%1f%(refname)%1f%(refname:short)%1f%(objectname:short)%1f%(committerdate:format:%Y-%m-%d %H:%M)%1f%(subject)",
		"refs/heads",
		"refs/remotes",
	)
	if err != nil {
		message := strings.TrimSpace(string(out))
		if message == "" {
			message = err.Error()
		}
		_, _ = errorColor.Printf("❌ 分支获取失败: %s\n", message)
		return nil
	}

	entries := parseEntries(string(out))
	if len(entries) == 0 {
		_, _ = infoColor.Println("未找到分支")
		return nil
	}

	_, _ = titleColor.Println("🌿 分支列表")
	_, _ = infoColor.Println()

	nameWidth := branchNameWidth(entries)
	for _, entry := range entries {
		printEntry(entry, nameWidth)
	}

	return nil
}

func parseEntries(out string) []Entry {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	entries := make([]Entry, 0, len(lines))

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.SplitN(line, "\x1f", 6)
		if len(parts) != 6 {
			continue
		}

		fullName := parts[1]
		name := parts[2]
		if strings.HasSuffix(name, "/HEAD") {
			continue
		}

		entries = append(entries, Entry{
			Name:      name,
			Kind:      branchKind(fullName),
			Hash:      parts[3],
			Date:      parts[4],
			Subject:   parts[5],
			IsCurrent: strings.TrimSpace(parts[0]) == "*",
		})
	}

	return entries
}

func branchKind(name string) string {
	if strings.HasPrefix(name, "refs/remotes/") {
		return "远程"
	}
	return "本地"
}

func branchNameWidth(entries []Entry) int {
	width := 32
	for _, entry := range entries {
		if w := term.DisplayWidth(entry.Name); w > width {
			width = w
		}
	}
	return width
}

func printEntry(entry Entry, nameWidth int) {
	marker := " "
	if entry.IsCurrent {
		marker = "*"
	}
	_, _ = currentColor.Printf("%s ", marker)

	if entry.Kind == "本地" {
		_, _ = localColor.Printf("%s ", term.PadRight(entry.Kind, 6))
	} else {
		_, _ = remoteColor.Printf("%s ", term.PadRight(entry.Kind, 6))
	}

	_, _ = infoColor.Printf("%s ", term.PadRight(entry.Name, nameWidth))
	_, _ = hashColor.Printf("%-10s ", entry.Hash)
	_, _ = dateColor.Printf("%-16s ", entry.Date)
	_, _ = infoColor.Println(entry.Subject)
}
