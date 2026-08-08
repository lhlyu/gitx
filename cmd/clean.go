package cmd

import (
	"github.com/lhlyu/gitx/internal/clean"
	"github.com/spf13/cobra"
)

func init() {
	cleanCmd.Flags().BoolVarP(&cleanDryRun, "dry-run", "n", false, "仅显示将执行的操作，不修改仓库")
	rootCmd.AddCommand(cleanCmd)
}

var cleanDryRun bool

var cleanCmd = &cobra.Command{
	Use:   "clean [depth]",
	Short: "清理仓库，重置到最新提交状态",
	Long:  "清理仓库，重置到最新提交状态，depth 参数指定扫描深度（默认为 0，表示只清理当前目录）；使用 -n 可仅预览而不修改仓库",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		depth, err := parseOptionalInt(args, 0, 0, "depth")
		if err != nil {
			return err
		}
		return clean.Run(depth, cleanDryRun)
	},
}
