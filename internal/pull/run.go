package pull

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
	results := repo.ProcessWithProgress(targets, "拉取中", func(t repo.Target) Result {
		return pullRepo(client, t)
	})

	_, _ = titleColor.Println("🔄 拉取代码结果")
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

func pullRepo(client *git.Client, t repo.Target) Result {
	out, err := client.RunInDir(t.Path, "pull")
	if err != nil {
		return Result{
			Name:    t.Name,
			Success: false,
			Message: strings.TrimSpace(string(out)),
		}
	}

	return Result{
		Name:    t.Name,
		Success: true,
		Message: classifyPullOutput(string(out)),
	}
}

func classifyPullOutput(out string) string {
	output := strings.ToLower(strings.TrimSpace(out))
	if strings.Contains(output, "already up to date") || strings.Contains(output, "already up-to-date") {
		return "已是最新"
	}
	return "已更新"
}
