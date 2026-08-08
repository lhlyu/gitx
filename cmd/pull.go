package cmd

import (
	"github.com/lhlyu/gitx/internal/pull"
	"github.com/spf13/cobra"
)

func init() {
	pullCmd.Flags().BoolVar(&pullRebase, "rebase", false, "将本地提交重放到远程最新提交之后，保持线性历史")
	pullCmd.Flags().BoolVar(&pullFFOnly, "ff-only", false, "仅允许快进；本地与远程分叉时直接失败")
	pullCmd.Flags().BoolVar(&pullPrune, "prune", false, "清理远程已删除的远程跟踪分支，不删除本地分支")
	pullCmd.MarkFlagsMutuallyExclusive("rebase", "ff-only")
	rootCmd.AddCommand(pullCmd)
}

var (
	pullRebase bool
	pullFFOnly bool
	pullPrune  bool
)

var pullCmd = &cobra.Command{
	Use:     "pull [depth]",
	Aliases: []string{"p"},
	Short:   "拉取最新代码",
	Long:    "拉取最新代码，depth 参数指定扫描深度；未传 depth 时，当前目录是 Git 仓库则只拉取当前仓库，否则扫描直接子目录",
	Example: "  gitx pull --ff-only\n  gitx p 1 --rebase --prune",
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		depth, err := parseOptionalInt(args, pull.AutoDepth, 0, "depth")
		if err != nil {
			return err
		}
		return pull.Run(pull.Options{
			Depth:  depth,
			Rebase: pullRebase,
			FFOnly: pullFFOnly,
			Prune:  pullPrune,
		})
	},
}
