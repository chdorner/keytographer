package cmd

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
)

func addConfigFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("config", "c", "", "path to the keymap configuration file")
}

func parseConfigFlag(cmd *cobra.Command) (string, error) {
	config, _ := cmd.Flags().GetString("config")
	if config == "" {
		return "", errors.New("missing path to keymap configuration file")
	}
	_, err := os.Stat(config)
	if err != nil {
		return "", errors.New("specified keymap configuration file does not exist")
	}
	return config, nil
}
