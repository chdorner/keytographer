package cmd

import (
	"context"
	"errors"
	"os"

	"github.com/chdorner/keytographer/internal/keytographer"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewValidateCommand() *cobra.Command {
	var ctx context.Context
	var configFile string

	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate keymap configuration.",

		PreRunE: func(cmd *cobra.Command, args []string) error {
			ctx = createContext(cmd.Flags())
			configureLogging(ctx)

			configFile, _ = cmd.Flags().GetString("config")
			if configFile == "" {
				return errors.New("missing path to keymap configuration file")
			}
			_, err := os.Stat(configFile)
			if err != nil {
				return errors.New("specified keymap configuration file does not exist")
			}

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := keytographer.Load(configFile)
			if err != nil {
				return err
			}

			err = keytographer.Validate(data)
			if err != nil {
				logrus.Error(err)
				os.Exit(1)
				return nil
			}

			logrus.Info("Configuration is valid!")
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("out", "o", "", "Path to output file.")

	return cmd
}
