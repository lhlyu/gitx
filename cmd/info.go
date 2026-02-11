package cmd

import (
	"github.com/lhlyu/gitx/internal/info"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(infoCmd)
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "显示仓库信息",
	RunE: func(cmd *cobra.Command, args []string) error {
		return info.Run()
	},
}
