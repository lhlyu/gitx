package status

import (
	"fmt"
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
		_, _ = projectColor.Printf("%s ", term.PadRight(r.Name, 40))
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
	_, _ = infoColor.Printf("%s ", term.PadRight(branch, 18))

	// 同步状态：领先/落后远程
	switch {
	case s.NoUpstream:
		_, _ = infoColor.Printf("%s ", term.PadRight("无上游", 12))
	case s.Ahead > 0 && s.Behind > 0:
		ahead := fmt.Sprintf("↑%d", s.Ahead)
		behind := fmt.Sprintf("↓%d", s.Behind)
		_, _ = aheadColor.Print(ahead)
		_, _ = behindColor.Print(behind)
		_, _ = infoColor.Printf("%s ", padding(term.DisplayWidth(ahead)+term.DisplayWidth(behind), 12))
	case s.Ahead > 0:
		_, _ = aheadColor.Printf("%s ", term.PadRight(fmt.Sprintf("↑%d", s.Ahead), 12))
	case s.Behind > 0:
		_, _ = behindColor.Printf("%s ", term.PadRight(fmt.Sprintf("↓%d", s.Behind), 12))
	default:
		_, _ = infoColor.Printf("%s ", term.PadRight("已同步", 12))
	}

	// 工作区状态
	if s.IsClean() {
		_, _ = cleanColor.Println("✅ 干净")
	} else {
		_, _ = dirtyColor.Printf("❌ %d 个改动\n", s.ChangedFiles)
	}
}

func padding(currentWidth, targetWidth int) string {
	if n := targetWidth - currentWidth; n > 0 {
		return strings.Repeat(" ", n)
	}
	return ""
}
