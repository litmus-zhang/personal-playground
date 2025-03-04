package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// docsCmd represents the docs command
var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Generate documentation for the application",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := cmd.Flags().GetString("dir")
		if err != nil {
			return err
		}

		if dir == "" {
			dir = os.TempDir()
			if err := os.MkdirAll(dir, 0755); err != nil {
				return err
			}
		}
		return docsAction(os.Stdout, dir)
	},
}

func docsAction(out io.Writer, dir string) error {
	if err := doc.GenMarkdownTree(rootCmd, dir); err != nil {
		return err
	}
	_, err := fmt.Fprintf(out, "Documentation generated in %s\n", dir)
	return err
}

func init() {
	rootCmd.AddCommand(docsCmd)

	docsCmd.Flags().StringP("dir", "d", "", "Destination directory for docs")
}
