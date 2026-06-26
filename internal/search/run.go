package search

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/lhlyu/gitx/internal/git"
)

var (
	titleColor = color.New(color.FgCyan, color.Bold)
	matchColor = color.New(color.FgYellow)
	infoColor  = color.New(color.FgWhite)
	errorColor = color.New(color.FgRed, color.Bold)
)

func Run(keyword string) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	root, err := repoRoot(currentDir)
	if err != nil {
		_, _ = errorColor.Println("❌ 当前目录不是 Git 项目")
		return nil
	}

	out, err := runRipgrep(root, keyword)
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
			_, _ = infoColor.Printf("未找到关键字: %s\n", keyword)
			return nil
		}

		message := strings.TrimSpace(string(out))
		if message == "" {
			message = err.Error()
		}
		return fmt.Errorf("搜索失败: %s", message)
	}

	matches := parseMatches(string(out))
	if len(matches) == 0 {
		_, _ = infoColor.Printf("未找到关键字: %s\n", keyword)
		return nil
	}

	_, _ = titleColor.Printf("🔎 找到 %d 条匹配\n", len(matches))
	_, _ = infoColor.Println()
	for _, match := range matches {
		_, _ = matchColor.Println(match)
	}

	return nil
}

func repoRoot(dir string) (string, error) {
	client := git.NewClient()
	out, err := client.RunInDir(dir, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", err
	}
	root := strings.TrimSpace(string(out))
	if root == "" {
		return "", fmt.Errorf("empty git root")
	}
	return root, nil
}

func runRipgrep(root, keyword string) ([]byte, error) {
	cmd := exec.Command(
		"rg",
		"--line-number",
		"--with-filename",
		"--fixed-strings",
		"--hidden",
		"--glob", "!.git/**",
		"--",
		keyword,
	)
	cmd.Dir = root
	return cmd.CombinedOutput()
}

func parseMatches(out string) []string {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	matches := make([]string, 0, len(lines))

	for _, line := range lines {
		match := strings.TrimSpace(line)
		if match == "" {
			continue
		}
		matches = append(matches, match)
	}

	return matches
}
