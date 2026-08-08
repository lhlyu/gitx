package cmd

import (
	"github.com/lhlyu/gitx/internal/undo"
	"github.com/spf13/cobra"
)

func init() {
	undoCmd.Flags().BoolVarP(&undoDryRun, "dry-run", "n", false, "仅显示将执行的操作，不修改仓库")
	rootCmd.AddCommand(undoCmd)
}

var undoDryRun bool

var undoCmd = &cobra.Command{
	Use:   "undo [depth]",
	Short: "撤销工作区和暂存区的修改",
	Long:  "撤销工作区和暂存区的修改，depth 参数指定扫描深度（默认为 0，表示只撤销当前目录）；使用 -n 可仅预览而不修改仓库",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		depth, err := parseOptionalInt(args, 0, 0, "depth")
		if err != nil {
			return err
		}
		return undo.Run(depth, undoDryRun)
	},
}
