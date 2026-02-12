package list

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/lhlyu/gitx/internal/git"
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
	Path    string
	IsClean bool
}

func Run(depth int) error {
	if depth < 1 {
		depth = 1
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	projects, err := scanProjects(currentDir, depth, 0)
	if err != nil {
		return err
	}

	if len(projects) == 0 {
		_, _ = infoColor.Println("未找到 Git 项目")
		return nil
	}

	_, _ = titleColor.Println("📁 Git 项目列表")
	_, _ = infoColor.Println()

	for _, proj := range projects {
		if proj.IsClean {
			_, _ = projectColor.Printf("%-50s ", proj.Name)
			_, _ = cleanColor.Println("✅")
		} else {
			_, _ = projectColor.Printf("%-50s ", proj.Name)
			_, _ = dirtyColor.Println("❌")
		}
	}

	return nil
}

func scanProjects(dir string, maxDepth, currentDepth int) ([]*Project, error) {
	var projects []*Project

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		projectPath := filepath.Join(dir, entry.Name())
		gitPath := filepath.Join(projectPath, ".git")

		// 检查是否是 Git 项目
		if _, err := os.Stat(gitPath); err == nil || os.IsExist(err) {
			isClean := checkIfClean(projectPath)
			// 计算相对路径显示
			relPath, _ := filepath.Rel(dir, projectPath)
			if currentDepth > 0 {
				relPath = filepath.Join(strings.Repeat("../", currentDepth), relPath)
			}
			projects = append(projects, &Project{
				Name:    relPath,
				Path:    projectPath,
				IsClean: isClean,
			})
		} else if currentDepth < maxDepth-1 {
			// 继续递归扫描子目录
			subProjects, err := scanProjects(projectPath, maxDepth, currentDepth+1)
			if err == nil {
				projects = append(projects, subProjects...)
			}
		}
	}

	return projects, nil
}

func checkIfClean(projectPath string) bool {
	client := git.NewClient()

	// 临时切换到项目目录执行 git 命令
	originalDir, _ := os.Getwd()
	defer func(dir string) {
		_ = os.Chdir(dir)
	}(originalDir)

	if err := os.Chdir(projectPath); err != nil {
		return false
	}

	out, err := client.Run("status", "--porcelain")
	if err != nil {
		return false
	}

	return strings.TrimSpace(string(out)) == ""
}
