package cmd

import (
	"fmt"

	"github.com/lhlyu/gitx/internal/search"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(searchCmd)
}

var searchCmd = &cobra.Command{
	Use:   "search <keyword>",
	Short: "快速查找关键字出现在哪些文件里",
	Long:  "快速查找关键字出现在哪些文件里，按字面量匹配，不把关键字当正则表达式解析",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("需要提供 1 个关键字")
		}
		if args[0] == "" {
			return fmt.Errorf("关键字不能为空")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return search.Run(args[0])
	},
}
