package cmd

import (
	gittag "github.com/lhlyu/gitx/internal/tag"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(tagCmd)
}

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "查看当前仓库的本地 tag",
	Long:  "查看当前仓库的本地 tag，按时间倒序显示类型、时间、目标提交和提交标题",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return gittag.Run()
	},
}
