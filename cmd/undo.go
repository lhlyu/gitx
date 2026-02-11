package cmd

import (
	"github.com/lhlyu/gitx/internal/undo"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(undoCmd)
}

var undoCmd = &cobra.Command{
	Use:   "undo",
	Short: "撤销工作区和暂存区的修改",
	RunE: func(cmd *cobra.Command, args []string) error {
		return undo.Run()
	},
}
