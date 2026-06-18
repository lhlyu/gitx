package cmd

import (
	"strconv"

	"github.com/lhlyu/gitx/internal/status"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   "status [depth]",
	Short: "查看仓库与远程的同步状态",
	Long:  "查看仓库分支、领先/落后远程的提交数及工作区状态（基于上次 fetch 的远程信息，默认 depth=0 只看当前目录）",
	RunE: func(cmd *cobra.Command, args []string) error {
		depth := 0
		if len(args) > 0 {
			if d, err := strconv.Atoi(args[0]); err == nil && d >= 0 {
				depth = d
			}
		}
		return status.Run(depth)
	},
}
