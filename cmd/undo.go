package cmd

import (
	"strconv"

	"github.com/lhlyu/gitx/internal/undo"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(undoCmd)
}

var undoCmd = &cobra.Command{
	Use:   "undo [depth]",
	Short: "撤销工作区和暂存区的修改",
	Long:  "撤销工作区和暂存区的修改，depth 参数指定扫描深度（默认为 0，表示只撤销当前目录）",
	RunE: func(cmd *cobra.Command, args []string) error {
		depth := 0
		if len(args) > 0 {
			if d, err := strconv.Atoi(args[0]); err == nil && d >= 0 {
				depth = d
			}
		}
		return undo.Run(depth)
	},
}
