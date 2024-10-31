package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/litmus-zhang/pScan/scan"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:          "add <host1> <host2> ...",
	Aliases:      []string{"a"},
	Short:        "Add hosts to the hosts list",
	Args:         cobra.MinimumNArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		hostsFile := viper.GetString("hosts-file")

		return addAction(os.Stdout, hostsFile, args)
	},
}

func addAction(out io.Writer, hostsFile string, hosts []string) error {
	hl := &scan.HostsList{}
	if err := hl.Load(hostsFile); err != nil {
		return err
	}

	for _, h := range hosts {
		if err := hl.Add(h); err != nil {
			return err
		}
		fmt.Fprintf(out, "Added hosts:  %s\n", h)
	}
	return hl.Save(hostsFile)

}

func init() {
	hostsCmd.AddCommand(addCmd)
}
