package cmd

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(upgradeCmd)
}

var upgradeCmd = &cobra.Command{
	Use:     "upgrade",
	Aliases: []string{"up"},
	Short:   "更新 gitx 到最新版本",
	Long:    "通过 go install github.com/lhlyu/gitx@latest 更新 gitx 到最新版本",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		upgrade := exec.Command("go", "install", "github.com/lhlyu/gitx@latest")
		upgrade.Stdin = os.Stdin
		upgrade.Stdout = os.Stdout
		upgrade.Stderr = os.Stderr
		return upgrade.Run()
	},
}
