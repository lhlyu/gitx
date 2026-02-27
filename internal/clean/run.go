package clean

import (
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/lhlyu/gitx/internal/git"
)

var (
	titleColor   = color.New(color.FgCyan, color.Bold)
	projectColor = color.New(color.FgYellow)
	successColor = color.New(color.FgGreen, color.Bold)
	errorColor   = color.New(color.FgRed, color.Bold)
	infoColor    = color.New(color.FgWhite)
	warningColor = color.New(color.FgYellow, color.Bold)
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

	_, _ = warningColor.Println("⚠️  警告：此操作将清除所有未提交的修改和未跟踪的文件")

	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	var results []*Result

	if depth == 0 {
		// 只清理当前目录
		if isGitRepo(currentDir) {
			result := cleanRepo(currentDir, filepath.Base(currentDir))
			results = append(results, result)
		} else {
			_, _ = errorColor.Println("❌ 当前目录不是 Git 项目")
			return nil
		}
	} else {
		// 清理指定深度的所有 Git 项目
		results = scanAndClean(currentDir, depth, 0)
	}

	if len(results) == 0 {
		_, _ = infoColor.Println("未找到 Git 项目")
		return nil
	}

	_, _ = titleColor.Println("🧹 清理结果")
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

func scanAndClean(dir string, maxDepth, currentDepth int) []*Result {
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
			result := cleanRepo(projectPath, entry.Name())
			results = append(results, result)
		} else if currentDepth < maxDepth-1 {
			// 继续递归扫描子目录
			subResults := scanAndClean(projectPath, maxDepth, currentDepth+1)
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

func cleanRepo(projectPath, projectName string) *Result {
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

	// 执行 git reset --hard HEAD
	if _, err := client.Run("reset", "--hard", "HEAD"); err != nil {
		return &Result{
			Name:    projectName,
			Path:    projectPath,
			Success: false,
			Message: "重置失败",
		}
	}

	// 执行 git clean -fd
	if _, err := client.Run("clean", "-fd"); err != nil {
		return &Result{
			Name:    projectName,
			Path:    projectPath,
			Success: false,
			Message: "清理失败",
		}
	}

	return &Result{
		Name:    projectName,
		Path:    projectPath,
		Success: true,
		Message: "清理成功",
	}
}
