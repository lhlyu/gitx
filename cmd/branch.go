package cmd

import (
	branch "github.com/lhlyu/gitx/internal/branch"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(branchCmd)
}

var branchCmd = &cobra.Command{
	Use:   "branch",
	Short: "查看当前仓库的本地与远程分支",
	Long:  "查看当前仓库已知的本地分支和远程跟踪分支（基于本地 refs，不主动联网 fetch）",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return branch.Run()
	},
}
