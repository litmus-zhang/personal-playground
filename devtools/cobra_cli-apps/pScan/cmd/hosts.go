package cmd

import (
	"github.com/spf13/cobra"
)

// hostsCmd represents the hosts command
var hostsCmd = &cobra.Command{
	Use:   "hosts",
	Short: "Manage the hosts list",
	Long:  `Manages the hosts lists for pScan.`,
}

func init() {
	rootCmd.AddCommand(hostsCmd)
}
