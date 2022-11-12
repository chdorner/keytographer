package cmd

import (
	"context"

	"github.com/chdorner/keytographer/live"
	"github.com/chdorner/keytographer/renderer"
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

			var err error
			configFile, err = parseConfigFlag(cmd)
			if err != nil {
				return err
			}

			host, _ = cmd.Flags().GetString("host")
			port, _ = cmd.Flags().GetInt("port")

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			renderer := renderer.NewRenderer()
			server, err := live.NewServer(ctx, renderer, configFile, host, port)
			if err != nil {
				return err
			}
			logrus.Infof("starting server on %s:%d", host, port)
			return server.ListenAndServe()
		},
	}

	fl := cmd.Flags()
	addConfigFlag(cmd)
	fl.StringP("host", "H", "localhost", "Host on which to run the live server on.")
	fl.IntP("port", "p", 8080, "Port on which to run the live server on.")

	return cmd
}
