package clean

import (
	"os"

	"github.com/fatih/color"
	"github.com/lhlyu/gitx/internal/git"
	"github.com/lhlyu/gitx/internal/repo"
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
	Success bool
	Message string
}

func Run(depth int) error {
	_, _ = warningColor.Println("⚠️  警告：此操作将清除所有未提交的修改和未跟踪的文件")

	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	if depth == 0 && !repo.IsGitRepo(currentDir) {
		_, _ = errorColor.Println("❌ 当前目录不是 Git 项目")
		return nil
	}

	targets := repo.Scan(currentDir, depth)
	if len(targets) == 0 {
		_, _ = infoColor.Println("未找到 Git 项目")
		return nil
	}

	client := git.NewClient()
	results := repo.ProcessWithProgress(targets, "清理中", func(t repo.Target) Result {
		return cleanRepo(client, t)
	})

	_, _ = titleColor.Println("🧹 清理结果")
	_, _ = infoColor.Println()

	for _, result := range results {
		_, _ = projectColor.Printf("%-50s ", result.Name)
		if result.Success {
			_, _ = successColor.Printf("✅ %s\n", result.Message)
		} else {
			_, _ = errorColor.Printf("❌ %s\n", result.Message)
		}
	}

	return nil
}

func cleanRepo(client *git.Client, t repo.Target) Result {
	if out, err := client.RunInDir(t.Path, "reset", "--hard", "HEAD"); err != nil {
		return Result{Name: t.Name, Success: false, Message: "重置失败: " + repo.FirstLine(out)}
	}

	if out, err := client.RunInDir(t.Path, "clean", "-fd"); err != nil {
		return Result{Name: t.Name, Success: false, Message: "清理失败: " + repo.FirstLine(out)}
	}

	return Result{Name: t.Name, Success: true, Message: "清理成功"}
}
