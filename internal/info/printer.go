package info

import (
	"fmt"

	"github.com/fatih/color"
)

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
	_, _ = labelColor.Printf("%-14s : %s\n", label, value)
}

func printWorking(i *Info) {
	_, _ = labelColor.Printf("%-14s : ", "工作区")

	if i.IsClean {
		_, _ = color.New(color.FgGreen, color.Bold).Println("干净 ✅")
	} else {
		_, _ = color.New(color.FgRed, color.Bold).Printf("有改动 ❌（%d 个已修改文件）\n", i.ChangedFiles)
	}
}
