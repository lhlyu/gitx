package cmd

import (
	"strconv"

	"github.com/lhlyu/gitx/internal/list"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list [depth]",
	Short: "列出当前目录下的 Git 项目",
	Long:  "列出当前目录下的 Git 项目，depth 参数指定扫描深度（默认为 1）",
	RunE: func(cmd *cobra.Command, args []string) error {
		depth := 1
		if len(args) > 0 {
			if d, err := strconv.Atoi(args[0]); err == nil && d > 0 {
				depth = d
			}
		}
		return list.Run(depth)
	},
}
