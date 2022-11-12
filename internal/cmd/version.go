package cmd

import (
	"fmt"

	"github.com/chdorner/keytographer/internal/keytographer"
	"github.com/spf13/cobra"
)

func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version of keytographer.",

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("keytographer version %s\n", keytographer.Version)
		},
	}
}
