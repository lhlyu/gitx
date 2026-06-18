package list

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
	cleanColor   = color.New(color.FgGreen, color.Bold)
	dirtyColor   = color.New(color.FgRed, color.Bold)
	infoColor    = color.New(color.FgWhite)
)

type Project struct {
	Name    string
	IsClean bool
	Branch  string
}

func Run(depth int) error {
	if depth < 1 {
		depth = 1
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	targets := repo.Scan(currentDir, depth)
	if len(targets) == 0 {
		_, _ = infoColor.Println("未找到 Git 项目")
		return nil
	}

	client := git.NewClient()
	projects := repo.Process(targets, func(t repo.Target) Project {
		return inspect(client, t)
	})

	_, _ = titleColor.Println("📁 Git 项目列表")
	_, _ = infoColor.Println()

	for _, proj := range projects {
		branch := proj.Branch
		if branch == "" {
			branch = "(unknown)"
		}
		_, _ = projectColor.Printf("%-50s ", proj.Name)
		_, _ = infoColor.Printf("%-18s ", branch)
		if proj.IsClean {
			_, _ = cleanColor.Println("✅")
		} else {
			_, _ = dirtyColor.Println("❌")
		}
	}

	return nil
}

// inspect 用一次 `git status -b --porcelain` 同时取分支与工作区状态。
func inspect(client *git.Client, t repo.Target) Project {
	proj := Project{Name: t.Name}

	out, err := client.RunInDir(t.Path, "status", "-b", "--porcelain")
	if err != nil {
		return proj
	}

	lines := strings.Split(strings.TrimRight(string(out), "\n"), "\n")
	proj.IsClean = true
	for i, line := range lines {
		if i == 0 && strings.HasPrefix(line, "## ") {
			proj.Branch = parseBranch(line)
			continue
		}
		if strings.TrimSpace(line) != "" {
			proj.IsClean = false
		}
	}

	return proj
}

// parseBranch 解析 `## main...origin/main [ahead 1]` 这类分支行，
// 同时兼容 `## No commits yet on main` 与 `## HEAD (no branch)`。
func parseBranch(line string) string {
	branch := strings.TrimPrefix(line, "## ")

	if strings.HasPrefix(branch, "HEAD (no branch)") {
		return "(detached)"
	}
	if rest, ok := strings.CutPrefix(branch, "No commits yet on "); ok {
		return strings.TrimSpace(rest)
	}

	if idx := strings.Index(branch, "..."); idx >= 0 {
		branch = branch[:idx]
	}
	if idx := strings.Index(branch, " "); idx >= 0 {
		branch = branch[:idx]
	}
	return branch
}
