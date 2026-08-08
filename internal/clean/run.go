package clean

import (
	"os"

	"github.com/fatih/color"
	"github.com/lhlyu/gitx/internal/git"
	"github.com/lhlyu/gitx/internal/repo"
	"github.com/lhlyu/gitx/internal/term"
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

func Run(depth int, dryRun bool) error {
	if dryRun {
		_, _ = warningColor.Println("预览模式：不会修改仓库")
	} else {
		_, _ = warningColor.Println("⚠️  警告：此操作将清除所有未提交的修改和未跟踪的文件")
	}

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
	label := "清理中"
	title := "🧹 清理结果"
	if dryRun {
		label = "预览中"
		title = "清理预览"
	}
	results := repo.ProcessWithProgress(targets, label, func(t repo.Target) Result {
		return cleanRepo(client, t, dryRun)
	})

	_, _ = titleColor.Println(title)
	_, _ = infoColor.Println()

	for _, result := range results {
		_, _ = projectColor.Printf("%s ", term.PadRight(result.Name, 50))
		if result.Success {
			_, _ = successColor.Printf("✅ %s\n", result.Message)
		} else {
			_, _ = errorColor.Printf("❌ %s\n", result.Message)
		}
	}

	return nil
}

func cleanRepo(client *git.Client, t repo.Target, dryRun bool) Result {
	if dryRun {
		return Result{Name: t.Name, Success: true, Message: "将执行: git reset --hard HEAD && git clean -fd"}
	}

	if out, err := client.RunInDir(t.Path, "reset", "--hard", "HEAD"); err != nil {
		return Result{Name: t.Name, Success: false, Message: "重置失败: " + repo.FirstLine(out)}
	}

	if out, err := client.RunInDir(t.Path, "clean", "-fd"); err != nil {
		return Result{Name: t.Name, Success: false, Message: "清理失败: " + repo.FirstLine(out)}
	}

	return Result{Name: t.Name, Success: true, Message: "清理成功"}
}
