package cmd

import (
	"fmt"
	"strconv"

	gitreset "github.com/lhlyu/gitx/internal/reset"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(resetCmd)
}

var resetCmd = &cobra.Command{
	Use:   "reset <steps>",
	Short: "重置当前仓库到前 N 个提交",
	Long:  "执行 git reset --hard HEAD~<steps>，steps 必须是非负整数",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		steps, err := strconv.Atoi(args[0])
		if err != nil || steps < 0 {
			return fmt.Errorf("steps 必须是非负整数")
		}

		return gitreset.Run(steps)
	},
}
