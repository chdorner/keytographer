package cmd

import (
	"context"
	"errors"
	"os"

	"github.com/chdorner/keytographer/config"
	"github.com/chdorner/keytographer/export"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewExportCommand() *cobra.Command {
	var ctx context.Context
	var configFile string
	var outFile string

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export a configuration to QMK keymaps.",

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

			outFile, _ = cmd.Flags().GetString("out")
			if outFile == "" {
				return errors.New("missing required out option")
			}

			return nil
		},

		Run: func(cmd *cobra.Command, args []string) {
			data, err := config.Load(configFile)
			if err != nil {
				logrus.WithField("error", err).Error("failed to load config file")
				os.Exit(1)
			}

			err = config.Validate(data)
			if err != nil {
				logrus.WithField("error", err).Error("configuration is invalid")
				os.Exit(1)
			}

			config, err := config.Parse(data)
			if err != nil {
				logrus.WithField("error", err).Error("failed to parse config")
				os.Exit(1)
			}

			exporter := export.NewCExporter(ctx)
			err = exporter.Export(config, outFile)
			if err != nil {
				logrus.WithField("error", err).Error("failed to export keymap")
				os.Exit(1)
			}
		},
	}

	fl := cmd.Flags()
	fl.StringP("out", "o", "", "Path to the output file.")

	return cmd
}
