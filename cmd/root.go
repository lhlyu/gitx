package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// version 版本号，每次修改提交时同步更新为当时的日期
const version = "2026.06.26-1"

var rootCmd = &cobra.Command{
	Use:     "gitx",
	Short:   "gitx - Git 命令行工具",
	Long:    "gitx 是一款 Git 工具 (version " + version + ")",
	Version: version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
