package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/litmus-zhang/pScan/scan"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:          "delete <host1> <host2> ... ",
	Aliases:      []string{"d"},
	Short:        "Delete hosts from the hosts list",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		hostFile := viper.GetString("hosts-file")

		return deleteAction(os.Stdout, hostFile, args)
	},
}

func deleteAction(out io.Writer, hostFile string, hosts []string) error {
	hl := &scan.HostsList{}
	if err := hl.Load(hostFile); err != nil {
		return err
	}
	for _, h := range hosts {
		if err := hl.Remove(h); err != nil {
			return err
		}
		fmt.Fprintln(out, "Deleted host:", h)
	}
	return hl.Save(hostFile)
}

func init() {
	hostsCmd.AddCommand(deleteCmd)

}
