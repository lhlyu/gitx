package cmd

import (
	"fmt"
	"unicode/utf8"

	"github.com/lhlyu/gitx/internal/search"
	"github.com/spf13/cobra"
)

func init() {
	searchCmd.Flags().BoolVar(&searchAll, "all", false, "解除关键字长度和结果数量限制")
	rootCmd.AddCommand(searchCmd)
}

var searchAll bool

var searchCmd = &cobra.Command{
	Use:     "search <keyword>",
	Aliases: []string{"s"},
	Short:   "快速查找关键字出现在哪些文件里",
	Long:    "快速查找关键字出现在哪些文件里，按字面量匹配，不把关键字当正则表达式解析。默认关键字至少 2 个字符，最多显示 20 条匹配；使用 --all 解除限制",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("需要提供 1 个关键字")
		}
		if args[0] == "" {
			return fmt.Errorf("关键字不能为空")
		}
		all, err := cmd.Flags().GetBool("all")
		if err != nil {
			return err
		}
		if !all && utf8.RuneCountInString(args[0]) < 2 {
			return fmt.Errorf("关键字至少需要 2 个字符；如确实需要搜索单字符，请使用 --all")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		all, err := cmd.Flags().GetBool("all")
		if err != nil {
			return err
		}
		return search.Run(args[0], search.Options{All: all})
	},
}
