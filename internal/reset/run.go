package reset

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/lhlyu/gitx/internal/git"
)

var (
	successColor = color.New(color.FgGreen, color.Bold)
	errorColor   = color.New(color.FgRed, color.Bold)
	infoColor    = color.New(color.FgWhite)
	warningColor = color.New(color.FgYellow, color.Bold)
)

func Run(steps int) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	if !isGitRepo(currentDir) {
		_, _ = errorColor.Println("❌ 当前目录不是 Git 项目")
		return nil
	}

	_, _ = warningColor.Printf("⚠️  即将执行: git reset --hard HEAD~%d\n", steps)

	client := git.NewClient()
	target := fmt.Sprintf("HEAD~%d", steps)
	out, err := client.RunInDir(currentDir, "reset", "--hard", target)
	if err != nil {
		_, _ = errorColor.Printf("❌ 重置失败: %s\n", string(out))
		return nil
	}

	_, _ = successColor.Printf("✅ 重置成功: %s\n", filepath.Base(currentDir))
	if len(out) > 0 {
		_, _ = infoColor.Print(string(out))
	}

	return nil
}

func isGitRepo(dir string) bool {
	gitPath := filepath.Join(dir, ".git")
	if _, err := os.Stat(gitPath); err == nil {
		return true
	}
	return false
}
