package cmd

import (
	"context"
	"errors"
	"os"

	"github.com/chdorner/keymap-render/internal/keymap"
	"github.com/chdorner/keymap-render/internal/live"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewLiveCommand() *cobra.Command {
	var ctx context.Context
	var configFile string
	var host string
	var port int

	cmd := &cobra.Command{
		Use:   "live",
		Short: "Start a live server for easier render configuration.",

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

			host, _ = cmd.Flags().GetString("host")
			port, _ = cmd.Flags().GetInt("port")

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			renderer := keymap.NewRenderer()
			server, err := live.NewServer(ctx, renderer, configFile, host, port)
			if err != nil {
				return err
			}
			logrus.Infof("starting server on %s:%d", host, port)
			return server.ListenAndServe()
		},
	}

	fl := cmd.Flags()
	fl.StringP("config", "c", "", "Path to the keymap configuration file to watch for changes.")
	fl.StringP("host", "H", "localhost", "Host on which to run the live server on.")
	fl.IntP("port", "p", 8080, "Port on which to run the live server on.")

	return cmd
}
