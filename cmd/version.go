package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version string = "0.0.0-dev"
var CommitSHA string = "dirty"
var BuildTimestamp string = "now"

func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version of keytographer.",

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("keytographer version %s\n", BuildVersion())
		},
	}
}

func BuildVersion() string {
	return fmt.Sprintf("%s (%s, %s)", Version, CommitSHA, BuildTimestamp)
}
