package cmd

import (
	"strconv"

	"github.com/lhlyu/gitx/internal/clean"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cleanCmd)
}

var cleanCmd = &cobra.Command{
	Use:   "clean [depth]",
	Short: "清理仓库，重置到最新提交状态",
	Long:  "清理仓库，重置到最新提交状态，depth 参数指定扫描深度（默认为 0，表示只清理当前目录）",
	RunE: func(cmd *cobra.Command, args []string) error {
		depth := 0
		if len(args) > 0 {
			if d, err := strconv.Atoi(args[0]); err == nil && d >= 0 {
				depth = d
			}
		}
		return clean.Run(depth)
	},
}
