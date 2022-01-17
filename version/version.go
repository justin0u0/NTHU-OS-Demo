package version

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	Version string
	Commit  string
)

func NewVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "version",
		Short:   "Get version of the CLI",
		Example: "demo version",
		Run: func(_ *cobra.Command, _ []string) {
			pterm.FgGreen.Println("NTHU-OS-DEMO CLI")
			pterm.FgCyan.Printf("%-10s%s\n", "Version: ", Version)
			pterm.FgCyan.Printf("%-10s%s\n", "Commit: ", Commit)
		},
	}

	return cmd
}
