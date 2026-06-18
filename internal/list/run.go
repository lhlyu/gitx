package list

import (
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
		st := repo.Inspect(client, t.Path)
		return Project{Name: t.Name, IsClean: st.Err == nil && st.IsClean(), Branch: st.Branch}
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
