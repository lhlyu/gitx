package cmd

import (
	"fmt"
	"strconv"

	gitreset "github.com/lhlyu/gitx/internal/reset"
	"github.com/spf13/cobra"
)

func init() {
	resetCmd.Flags().BoolVarP(&resetDryRun, "dry-run", "n", false, "仅显示将执行的操作，不修改仓库")
	rootCmd.AddCommand(resetCmd)
}

var resetDryRun bool

var resetCmd = &cobra.Command{
	Use:   "reset <steps>",
	Short: "重置当前仓库到前 N 个提交",
	Long:  "执行 git reset --hard HEAD~<steps>，steps 必须是非负整数；使用 -n 可仅预览而不修改仓库",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		steps, err := strconv.Atoi(args[0])
		if err != nil || steps < 0 {
			return fmt.Errorf("steps 必须是非负整数")
		}

		return gitreset.Run(steps, resetDryRun)
	},
}
