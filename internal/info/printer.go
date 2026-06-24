package info

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/lhlyu/gitx/internal/term"
)

const labelWidth = 8

var (
	titleColor = color.New(color.FgCyan, color.Bold)
	labelColor = color.New(color.FgYellow)
)

func Print(i *Info) {
	_, _ = titleColor.Println("📦 仓库信息")
	fmt.Println()

	printField("用户名", i.UserName)
	printField("用户邮箱", i.UserEmail)
	printField("分支", i.Branch)
	printField("远程地址", i.RemoteURL)
	printWorking(i)
}

func printField(label, value string) {
	_, _ = labelColor.Printf("%s : %s\n", paddedLabel(label), value)
}

func printWorking(i *Info) {
	_, _ = labelColor.Printf("%s : ", paddedLabel("工作区"))

	if i.IsClean {
		_, _ = color.New(color.FgGreen, color.Bold).Println("干净 ✅")
	} else {
		_, _ = color.New(color.FgRed, color.Bold).Printf("有改动 ❌（%d 个已修改文件）\n", i.ChangedFiles)
	}
}

func paddedLabel(label string) string {
	return term.PadRight(label, labelWidth)
}
