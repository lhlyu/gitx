package status

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/lhlyu/gitx/internal/git"
	"github.com/lhlyu/gitx/internal/repo"
)

var (
	titleColor   = color.New(color.FgCyan, color.Bold)
	projectColor = color.New(color.FgYellow)
	cleanColor   = color.New(color.FgGreen, color.Bold)
	dirtyColor   = color.New(color.FgRed, color.Bold)
	aheadColor   = color.New(color.FgGreen)
	behindColor  = color.New(color.FgRed)
	infoColor    = color.New(color.FgWhite)
	errorColor   = color.New(color.FgRed, color.Bold)
)

type Result struct {
	Name   string
	Status repo.Status
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
	results := repo.ProcessWithProgress(targets, "检查中", func(t repo.Target) Result {
		return Result{Name: t.Name, Status: repo.Inspect(client, t.Path)}
	})

	_, _ = titleColor.Println("📊 仓库状态")
	_, _ = infoColor.Println()

	for _, r := range results {
		_, _ = projectColor.Printf("%-40s ", r.Name)
		printStatus(r.Status)
	}

	return nil
}

func printStatus(s repo.Status) {
	if s.Err != nil {
		_, _ = errorColor.Println("❌ 状态获取失败")
		return
	}

	branch := s.Branch
	if branch == "" {
		branch = "(unknown)"
	}
	_, _ = infoColor.Printf("%-18s ", branch)

	// 同步状态：领先/落后远程
	switch {
	case s.NoUpstream:
		_, _ = infoColor.Printf("%-12s ", "无上游")
	case s.Ahead > 0 && s.Behind > 0:
		_, _ = aheadColor.Printf("↑%d", s.Ahead)
		_, _ = behindColor.Printf("↓%d", s.Behind)
		_, _ = infoColor.Printf("%-*s ", pad(s.Ahead, s.Behind), "")
	case s.Ahead > 0:
		_, _ = aheadColor.Printf("%-12s ", fmt.Sprintf("↑%d", s.Ahead))
	case s.Behind > 0:
		_, _ = behindColor.Printf("%-12s ", fmt.Sprintf("↓%d", s.Behind))
	default:
		_, _ = infoColor.Printf("%-12s ", "已同步")
	}

	// 工作区状态
	if s.IsClean() {
		_, _ = cleanColor.Println("✅ 干净")
	} else {
		_, _ = dirtyColor.Printf("❌ %d 个改动\n", s.ChangedFiles)
	}
}

// pad 计算 "↑x↓y" 之后补齐到 12 列所需的空格数。
func pad(ahead, behind int) int {
	w := 2 + len(fmt.Sprint(ahead)) + len(fmt.Sprint(behind))
	if n := 12 - w; n > 0 {
		return n
	}
	return 1
}
