package cmd

import (
	"github.com/lhlyu/gitx/internal/clean"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cleanCmd)
}

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "清理仓库，重置到最新提交状态",
	RunE: func(cmd *cobra.Command, args []string) error {
		return clean.Run()
	},
}
