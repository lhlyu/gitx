package tag

import (
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/lhlyu/gitx/internal/git"
	"github.com/lhlyu/gitx/internal/repo"
	"github.com/lhlyu/gitx/internal/term"
)

var (
	titleColor = color.New(color.FgCyan, color.Bold)
	nameColor  = color.New(color.FgYellow)
	typeColor  = color.New(color.FgGreen)
	hashColor  = color.New(color.FgYellow)
	dateColor  = color.New(color.FgWhite)
	infoColor  = color.New(color.FgWhite)
	errorColor = color.New(color.FgRed, color.Bold)
)

type Entry struct {
	Name    string
	Kind    string
	Date    string
	Hash    string
	Subject string
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
		"refs/tags",
		"--sort=-creatordate",
		"--format=%(refname:short)%1f%(objecttype)%1f%(creatordate:format:%Y-%m-%d %H:%M)%1f%(*objectname:short)%1f%(objectname:short)%1f%(*subject)%1f%(subject)",
	)
	if err != nil {
		message := strings.TrimSpace(string(out))
		if message == "" {
			message = err.Error()
		}
		_, _ = errorColor.Printf("❌ tag 获取失败: %s\n", message)
		return nil
	}

	entries := parseEntries(string(out))
	if len(entries) == 0 {
		_, _ = infoColor.Println("未找到 tag")
		return nil
	}

	_, _ = titleColor.Println("🏷️ tag 列表")
	_, _ = infoColor.Println()

	nameWidth := tagNameWidth(entries)
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

		parts := strings.SplitN(line, "\x1f", 7)
		if len(parts) != 7 {
			continue
		}

		entry := Entry{
			Name:    parts[0],
			Kind:    tagKind(parts[1]),
			Date:    parts[2],
			Hash:    firstNonEmpty(parts[3], parts[4]),
			Subject: firstNonEmpty(parts[5], parts[6]),
		}
		entries = append(entries, entry)
	}

	return entries
}

func tagKind(objectType string) string {
	if objectType == "tag" {
		return "标注"
	}
	return "轻量"
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func tagNameWidth(entries []Entry) int {
	width := 12
	for _, entry := range entries {
		if w := term.DisplayWidth(entry.Name); w > width {
			width = w
		}
	}
	return width
}

func printEntry(entry Entry, nameWidth int) {
	_, _ = nameColor.Printf("%s ", term.PadRight(entry.Name, nameWidth))
	_, _ = typeColor.Printf("%s ", term.PadRight(entry.Kind, 6))
	_, _ = dateColor.Printf("%-16s ", entry.Date)
	_, _ = hashColor.Printf("%-10s ", entry.Hash)
	_, _ = infoColor.Println(entry.Subject)
}
