package pull

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/lhlyu/gitx/internal/git"
)

var (
	titleColor   = color.New(color.FgCyan, color.Bold)
	projectColor = color.New(color.FgYellow)
	successColor = color.New(color.FgGreen, color.Bold)
	errorColor   = color.New(color.FgRed, color.Bold)
	infoColor    = color.New(color.FgWhite)
)

const pullWorkerCount = 6

type Result struct {
	Name    string
	Path    string
	Success bool
	Message string
}

type repoTarget struct {
	Index int
	Name  string
	Path  string
}

func Run(depth int) error {
	if depth < 0 {
		depth = 0
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	var targets []*repoTarget

	if depth == 0 {
		// 只拉取当前目录
		if isGitRepo(currentDir) {
			targets = append(targets, &repoTarget{
				Index: 0,
				Name:  filepath.Base(currentDir),
				Path:  currentDir,
			})
		} else {
			_, _ = errorColor.Println("❌ 当前目录不是 Git 项目")
			return nil
		}
	} else {
		// 拉取指定深度的所有 Git 项目
		targets = scanRepos(currentDir, depth, 0)
	}

	if len(targets) == 0 {
		_, _ = infoColor.Println("未找到 Git 项目")
		return nil
	}

	for i, target := range targets {
		target.Index = i
	}

	results := pullRepos(targets)

	_, _ = titleColor.Println("🔄 拉取代码结果")
	_, _ = infoColor.Println()

	for _, result := range results {
		if result.Success {
			_, _ = projectColor.Printf("%-50s ", result.Name)
			_, _ = successColor.Printf("✅ %s\n", result.Message)
		} else {
			_, _ = projectColor.Printf("%-50s ", result.Name)
			_, _ = errorColor.Printf("❌ %s\n", result.Message)
		}
	}

	return nil
}

func scanRepos(dir string, maxDepth, currentDepth int) []*repoTarget {
	var targets []*repoTarget

	entries, err := os.ReadDir(dir)
	if err != nil {
		return targets
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		projectPath := filepath.Join(dir, entry.Name())

		if isGitRepo(projectPath) {
			targets = append(targets, &repoTarget{
				Name: entry.Name(),
				Path: projectPath,
			})
		} else if currentDepth < maxDepth-1 {
			// 继续递归扫描子目录
			subTargets := scanRepos(projectPath, maxDepth, currentDepth+1)
			targets = append(targets, subTargets...)
		}
	}

	return targets
}

func isGitRepo(dir string) bool {
	gitPath := filepath.Join(dir, ".git")
	if _, err := os.Stat(gitPath); err == nil {
		return true
	}
	return false
}

func pullRepos(targets []*repoTarget) []*Result {
	results := make([]*Result, len(targets))
	workerCount := pullWorkerCount
	if len(targets) < workerCount {
		workerCount = len(targets)
	}

	jobs := make(chan *repoTarget)
	var wg sync.WaitGroup

	for range workerCount {
		wg.Add(1)
		go func() {
			defer wg.Done()

			client := git.NewClient()
			for target := range jobs {
				results[target.Index] = pullRepo(client, target.Path, target.Name)
			}
		}()
	}

	for _, target := range targets {
		jobs <- target
	}
	close(jobs)
	wg.Wait()

	return results
}

func pullRepo(client *git.Client, projectPath, projectName string) *Result {
	out, err := client.RunInDir(projectPath, "pull")
	if err != nil {
		return &Result{
			Name:    projectName,
			Path:    projectPath,
			Success: false,
			Message: strings.TrimSpace(string(out)),
		}
	}

	return &Result{
		Name:    projectName,
		Path:    projectPath,
		Success: true,
		Message: classifyPullOutput(string(out)),
	}
}

func classifyPullOutput(out string) string {
	output := strings.ToLower(strings.TrimSpace(out))
	if output == "" {
		return "已更新"
	}
	if strings.Contains(output, "already up to date") || strings.Contains(output, "already up-to-date") {
		return "已是最新"
	}
	return "已更新"
}
