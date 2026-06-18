package undo

import (
	"os"
	"strings"

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
)

type Result struct {
	Name    string
	Success bool
	Message string
}

func Run(depth int) error {
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
	results := repo.Process(targets, func(t repo.Target) Result {
		return undoRepo(client, t)
	})

	_, _ = titleColor.Println("↩️  撤销结果")
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

func undoRepo(client *git.Client, t repo.Target) Result {
	out, err := client.RunInDir(t.Path, "status", "--porcelain")
	if err != nil {
		return Result{Name: t.Name, Success: false, Message: "状态获取失败"}
	}

	if strings.TrimSpace(string(out)) == "" {
		return Result{Name: t.Name, Success: true, Message: "无可撤销"}
	}

	// 一条命令原子地把暂存区和工作区都还原到 HEAD，避免分两步导致的中间状态。
	if _, err := client.RunInDir(t.Path, "restore", "--staged", "--worktree", "."); err != nil {
		return Result{Name: t.Name, Success: false, Message: "撤销失败"}
	}

	return Result{Name: t.Name, Success: true, Message: "撤销成功"}
}
