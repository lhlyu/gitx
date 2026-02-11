package undo

import (
	"github.com/fatih/color"
	"github.com/lhlyu/gitx/internal/git"
)

var (
	successColor = color.New(color.FgGreen, color.Bold)
	errorColor   = color.New(color.FgRed, color.Bold)
	infoColor    = color.New(color.FgCyan, color.Bold)
)

func Run() error {
	client := git.NewClient()

	_, _ = infoColor.Println("🔄 开始撤销本地修改...")

	// 执行 git restore .
	if _, err := client.Run("restore", "."); err != nil {
		_, _ = errorColor.Printf("❌ 撤销工作区修改失败: %v\n", err)
		return err
	}

	// 执行 git restore --staged .
	if _, err := client.Run("restore", "--staged", "."); err != nil {
		_, _ = errorColor.Printf("❌ 撤销暂存区修改失败: %v\n", err)
		return err
	}

	_, _ = successColor.Println("✅ 撤销成功！工作区和暂存区已恢复")
	return nil
}
