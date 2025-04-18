package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generate bash completion scripts",
	Long: `To load completion run
source <(pScan completion)

TO load completions automatically add to your .bashrc
# ~/.bashrc or ~/.profile
source <(pScan completion)
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return completionAction(os.Stdout)
	},
}

func completionAction(out io.Writer) error {
	return rootCmd.GenBashCompletion(out)
}

func init() {
	rootCmd.AddCommand(completionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// completionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// completionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
