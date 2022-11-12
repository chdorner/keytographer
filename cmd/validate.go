package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/chdorner/keytographer/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewValidateCommand() *cobra.Command {
	var ctx context.Context
	var configFile string

	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate keymap configuration",

		PreRunE: func(cmd *cobra.Command, args []string) error {
			ctx = createContext(cmd.Flags())
			configureLogging(ctx)

			var err error
			configFile, err = parseConfigFlag(cmd)
			if err != nil {
				return err
			}

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := config.Load(configFile)
			if err != nil {
				return err
			}

			err = config.Validate(data)
			if err != nil {
				logrus.Error(err)
				os.Exit(1)
				return nil
			}

			fmt.Println("Configuration is valid!")
			return nil
		},
	}

	addConfigFlag(cmd)

	return cmd
}
