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
)

var (
	titleColor = color.New(color.FgCyan, color.Bold)
	matchColor = color.New(color.FgYellow)
	infoColor  = color.New(color.FgWhite)
)

const defaultMatchLimit = 20

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
	matches   []string
	truncated bool
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
	for _, match := range result.matches {
		_, _ = matchColor.Println(match)
	}

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
		match := strings.TrimSpace(line)
		if match == "" {
			continue
		}
		if limit > 0 && len(result.matches) >= limit {
			result.truncated = true
			stop()
			break
		}
		result.matches = append(result.matches, match)
	}
	return result
}
