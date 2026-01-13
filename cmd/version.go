package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/EnSync-engine/CLI/pkg/version"
)

func newVersionCmd() *cobra.Command {
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			v := version.Get()

			if jsonOutput {
				return printJSON(cmd.OutOrStdout(), v)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "ensync %s\n", v.Version)
			if v.Commit != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "  commit: %s\n", v.Commit)
			}
			if v.BuildDate != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "  built:  %s\n", v.BuildDate)
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "output as JSON")

	return cmd
}
