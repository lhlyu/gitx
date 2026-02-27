package pull

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/lhlyu/gitx/internal/git"
)

var (
	titleColor   = color.New(color.FgCyan, color.Bold)
	projectColor = color.New(color.FgYellow)
	successColor = color.New(color.FgGreen, color.Bold)
	errorColor   = color.New(color.FgRed, color.Bold)
	infoColor    = color.New(color.FgWhite)
)

type Result struct {
	Name    string
	Path    string
	Success bool
	Message string
}

func Run(depth int) error {
	if depth < 0 {
		depth = 0
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	var results []*Result

	if depth == 0 {
		// 只拉取当前目录
		if isGitRepo(currentDir) {
			result := pullRepo(currentDir, filepath.Base(currentDir))
			results = append(results, result)
		} else {
			_, _ = errorColor.Println("❌ 当前目录不是 Git 项目")
			return nil
		}
	} else {
		// 拉取指定深度的所有 Git 项目
		results = scanAndPull(currentDir, depth, 0)
	}

	if len(results) == 0 {
		_, _ = infoColor.Println("未找到 Git 项目")
		return nil
	}

	_, _ = titleColor.Println("🔄 拉取代码结果")
	_, _ = infoColor.Println()

	for _, result := range results {
		if result.Success {
			_, _ = projectColor.Printf("%-50s ", result.Name)
			_, _ = successColor.Printf("✅ %s\n", result.Message)
		} else {
			_, _ = projectColor.Printf("%-50s ", result.Name)
			_, _ = errorColor.Printf("❌ %s\n", result.Message)
		}
	}

	return nil
}

func scanAndPull(dir string, maxDepth, currentDepth int) []*Result {
	var results []*Result

	entries, err := os.ReadDir(dir)
	if err != nil {
		return results
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		projectPath := filepath.Join(dir, entry.Name())

		if isGitRepo(projectPath) {
			result := pullRepo(projectPath, entry.Name())
			results = append(results, result)
		} else if currentDepth < maxDepth-1 {
			// 继续递归扫描子目录
			subResults := scanAndPull(projectPath, maxDepth, currentDepth+1)
			results = append(results, subResults...)
		}
	}

	return results
}

func isGitRepo(dir string) bool {
	gitPath := filepath.Join(dir, ".git")
	if _, err := os.Stat(gitPath); err == nil {
		return true
	}
	return false
}

func pullRepo(projectPath, projectName string) *Result {
	client := git.NewClient()

	// 临时切换到项目目录执行 git 命令
	originalDir, _ := os.Getwd()
	defer func(dir string) {
		_ = os.Chdir(dir)
	}(originalDir)

	if err := os.Chdir(projectPath); err != nil {
		return &Result{
			Name:    projectName,
			Path:    projectPath,
			Success: false,
			Message: "无法进入目录",
		}
	}

	out, err := client.Run("pull")
	if err != nil {
		return &Result{
			Name:    projectName,
			Path:    projectPath,
			Success: false,
			Message: strings.TrimSpace(string(out)),
		}
	}

	return &Result{
		Name:    projectName,
		Path:    projectPath,
		Success: true,
		Message: classifyPullOutput(string(out)),
	}
}

func classifyPullOutput(out string) string {
	output := strings.ToLower(strings.TrimSpace(out))
	if output == "" {
		return "已更新"
	}
	if strings.Contains(output, "already up to date") || strings.Contains(output, "already up-to-date") {
		return "已是最新"
	}
	return "已更新"
}
