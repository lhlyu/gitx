package clean

import (
	"github.com/fatih/color"
	"github.com/lhlyu/gitx/internal/git"
)

var (
	successColor = color.New(color.FgGreen, color.Bold)
	errorColor   = color.New(color.FgRed, color.Bold)
	infoColor    = color.New(color.FgCyan, color.Bold)
	warningColor = color.New(color.FgYellow, color.Bold)
)

func Run() error {
	client := git.NewClient()

	_, _ = warningColor.Println("⚠️  警告：此操作将清除所有未提交的修改和未跟踪的文件")
	_, _ = infoColor.Println("🧹 开始清理仓库...")

	// 执行 git reset --hard HEAD
	if _, err := client.Run("reset", "--hard", "HEAD"); err != nil {
		_, _ = errorColor.Printf("❌ 重置失败: %v\n", err)
		return err
	}

	// 执行 git clean -fd
	if _, err := client.Run("clean", "-fd"); err != nil {
		_, _ = errorColor.Printf("❌ 清理未跟踪文件失败: %v\n", err)
		return err
	}

	_, _ = successColor.Println("✅ 清理成功！仓库已恢复到最新提交状态")
	return nil
}
