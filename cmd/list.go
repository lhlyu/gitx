package cmd

import (
	"github.com/lhlyu/gitx/internal/list"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出当前目录下所有一级 Git 项目",
	RunE: func(cmd *cobra.Command, args []string) error {
		return list.Run()
	},
}
