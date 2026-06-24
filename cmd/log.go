package cmd

import (
	"fmt"
	"strconv"

	gitlog "github.com/lhlyu/gitx/internal/log"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(logCmd)
}

var logCmd = &cobra.Command{
	Use:   "log [n]",
	Short: "查看最近的提交记录",
	Long:  "查看最近的提交记录，n 表示最近多少条提交（默认 n=5）",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		n := 5
		if len(args) > 0 {
			parsed, err := strconv.Atoi(args[0])
			if err != nil || parsed <= 0 {
				return fmt.Errorf("n 必须是正整数")
			}
			n = parsed
		}

		return gitlog.Run(n)
	},
}
