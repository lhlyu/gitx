package search

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/lhlyu/gitx/internal/term"
)

var (
	titleColor = color.New(color.FgCyan, color.Bold)
	pathColor  = color.New(color.FgYellow)
	lineColor  = color.New(color.FgCyan)
	matchColor = color.New(color.FgYellow)
	infoColor  = color.New(color.FgWhite)
)

const defaultMatchLimit = 20
const rgFieldSeparator = "\x1f"
const pathColumnWidth = 50
const lineColumnWidth = 10

var defaultExcludeGlobs = []string{
	"!.git/**",
	"!node_modules/**",
	"!.idea/**",
	"!.vscode/**",
	"!.pnpm-store/**",
	"!.swc/**",
	"!.temp/**",
	"!.rn_temp/**",
	"!.cache/**",
	"!.nuxt/**",
	"!.output/**",
	"!.data/**",
	"!.nitro/**",
	"!.fleet/**",
	"!.DS_Store",
}

type Options struct {
	All bool
}

type searchResult struct {
	matches   []match
	truncated bool
}

type match struct {
	path string
	line string
	text string
}

func Run(keyword string, opts Options) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	limit := defaultMatchLimit
	if opts.All {
		limit = 0
	}

	result, err := runRipgrep(currentDir, keyword, limit)
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
			_, _ = infoColor.Printf("未找到关键字: %s\n", keyword)
			return nil
		}

		message := strings.TrimSpace(err.Error())
		if message == "" {
			message = err.Error()
		}
		return fmt.Errorf("搜索失败: %s", message)
	}

	if len(result.matches) == 0 {
		_, _ = infoColor.Printf("未找到关键字: %s\n", keyword)
		return nil
	}

	if result.truncated {
		_, _ = titleColor.Printf("🔎 找到至少 %d 条匹配，仅显示前 %d 条\n", len(result.matches)+1, len(result.matches))
		_, _ = infoColor.Println("使用 --all 查看全部匹配")
	} else {
		_, _ = titleColor.Printf("🔎 找到 %d 条匹配\n", len(result.matches))
	}
	_, _ = infoColor.Println()
	printMatches(result.matches, keyword)

	return nil
}

func runRipgrep(root, keyword string, limit int) (searchResult, error) {
	args := buildRipgrepArgs(keyword)
	cmd := exec.Command("rg", args...)
	cmd.Dir = root

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return searchResult{}, err
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return searchResult{}, err
	}

	result := collectMatches(stdout, limit, func() {
		_ = cmd.Process.Kill()
	})
	waitErr := cmd.Wait()
	if waitErr != nil && !result.truncated {
		message := strings.TrimSpace(stderr.String())
		if message != "" {
			return result, fmt.Errorf("%s", message)
		}
		return result, waitErr
	}

	return result, nil
}

func buildRipgrepArgs(keyword string) []string {
	args := []string{
		"--line-number",
		"--with-filename",
		"--field-match-separator",
		rgFieldSeparator,
		"--fixed-strings",
		"--hidden",
	}

	for _, glob := range defaultExcludeGlobs {
		args = append(args, "--glob", glob)
	}

	return append(args, "--", keyword)
}

func collectMatches(reader io.Reader, limit int, stop func()) searchResult {
	result := searchResult{}
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		if limit > 0 && len(result.matches) >= limit {
			result.truncated = true
			stop()
			break
		}
		result.matches = append(result.matches, parseMatchLine(line))
	}
	return result
}

func parseMatchLine(line string) match {
	parts := strings.SplitN(line, rgFieldSeparator, 3)
	if len(parts) != 3 {
		return match{text: line}
	}
	return match{
		path: parts[0],
		line: parts[1],
		text: parts[2],
	}
}

func printMatches(matches []match, keyword string) {
	for _, m := range matches {
		path := m.path
		if path == "" {
			path = "(无法解析文件)"
		}
		line := m.line
		if line != "" {
			line = "line:" + line
		}

		_, _ = pathColor.Printf("%s ", term.PadRight(path, pathColumnWidth))
		_, _ = lineColor.Printf("%s ", term.PadRight(line, lineColumnWidth))
		printHighlighted(m.text, keyword)
		_, _ = infoColor.Println()
	}
}

func printHighlighted(text, keyword string) {
	if keyword == "" {
		_, _ = infoColor.Print(text)
		return
	}

	for {
		index := strings.Index(text, keyword)
		if index < 0 {
			_, _ = infoColor.Print(text)
			return
		}
		_, _ = infoColor.Print(text[:index])
		_, _ = matchColor.Print(text[index : index+len(keyword)])
		text = text[index+len(keyword):]
	}
}
