package log

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/lhlyu/gitx/internal/git"
	"github.com/lhlyu/gitx/internal/repo"
	"github.com/lhlyu/gitx/internal/term"
)

var (
	titleColor  = color.New(color.FgCyan, color.Bold)
	hashColor   = color.New(color.FgYellow)
	dateColor   = color.New(color.FgWhite)
	authorColor = color.New(color.FgGreen)
	infoColor   = color.New(color.FgWhite)
	errorColor  = color.New(color.FgRed, color.Bold)
)

type Entry struct {
	Hash    string
	Date    string
	Author  string
	Subject string
}

func Run(n int) error {
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
		"log",
		"-n", fmt.Sprint(n),
		"--date=format:%Y-%m-%d %H:%M",
		"--pretty=format:%h%x1f%ad%x1f%an%x1f%s",
	)
	if err != nil {
		message := strings.TrimSpace(string(out))
		if message == "" {
			message = err.Error()
		}
		_, _ = errorColor.Printf("❌ 日志获取失败: %s\n", message)
		return nil
	}

	entries := parseEntries(string(out))
	if len(entries) == 0 {
		_, _ = infoColor.Println("暂无提交记录")
		return nil
	}

	_, _ = titleColor.Printf("🧾 最近 %d 次提交\n", n)
	_, _ = infoColor.Println()

	for _, entry := range entries {
		_, _ = hashColor.Printf("%-10s ", entry.Hash)
		_, _ = dateColor.Printf("%-16s ", entry.Date)
		_, _ = authorColor.Printf("%s ", term.PadRight(entry.Author, 18))
		_, _ = infoColor.Println(entry.Subject)
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

		parts := strings.SplitN(line, "\x1f", 4)
		if len(parts) != 4 {
			continue
		}

		entries = append(entries, Entry{
			Hash:    parts[0],
			Date:    parts[1],
			Author:  parts[2],
			Subject: parts[3],
		})
	}

	return entries
}
