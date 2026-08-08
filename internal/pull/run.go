package pull

import (
	"os"
	"strings"

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
)

// AutoDepth 根据当前目录是否为 Git 仓库自动选择扫描深度。
const AutoDepth = -1

type Result struct {
	Name    string
	Success bool
	Message string
}

type Options struct {
	Depth  int
	Rebase bool
	FFOnly bool
	Prune  bool
}

func Run(options Options) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	depth := resolveDepth(currentDir, options.Depth)

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
		return pullRepo(client, t, options)
	})

	_, _ = titleColor.Println("🔄 拉取代码结果")
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

func resolveDepth(currentDir string, depth int) int {
	if depth != AutoDepth {
		return depth
	}
	if repo.IsGitRepo(currentDir) {
		return 0
	}
	return 1
}

func pullRepo(client *git.Client, t repo.Target, options Options) Result {
	out, err := client.RunInDir(t.Path, pullArgs(options)...)
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

func pullArgs(options Options) []string {
	args := []string{"pull"}
	if options.Rebase {
		args = append(args, "--rebase")
	}
	if options.FFOnly {
		args = append(args, "--ff-only")
	}
	if options.Prune {
		args = append(args, "--prune")
	}
	return args
}

func classifyPullOutput(out string) string {
	output := strings.ToLower(strings.TrimSpace(out))
	if strings.Contains(output, "already up to date") || strings.Contains(output, "already up-to-date") {
		return "已是最新"
	}
	return "已更新"
}
