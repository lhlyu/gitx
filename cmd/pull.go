package cmd

import (
	"strconv"

	"github.com/lhlyu/gitx/internal/pull"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(pullCmd)
}

var pullCmd = &cobra.Command{
	Use:   "pull [depth]",
	Short: "拉取最新代码",
	Long:  "拉取最新代码，depth 参数指定扫描深度（默认为 0，表示只拉取当前目录）",
	RunE: func(cmd *cobra.Command, args []string) error {
		depth := 0
		if len(args) > 0 {
			if d, err := strconv.Atoi(args[0]); err == nil && d >= 0 {
				depth = d
			}
		}
		return pull.Run(depth)
	},
}
