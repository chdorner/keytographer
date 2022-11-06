package cmd

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var rootCmd = &cobra.Command{
	Use:   "keytographer",
	Short: "Render beautiful keymap visualizations.",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debug mode.")
	rootCmd.PersistentFlags().StringP("config", "c", "", "Path to the keymap configuration file to watch for changes.")

	rootCmd.AddCommand(NewValidateCommand())
	rootCmd.AddCommand(NewRenderCommand())
	rootCmd.AddCommand(NewLiveCommand())
}

func createContext(flags *pflag.FlagSet) context.Context {
	debug, _ := flags.GetBool("debug")
	return context.WithValue(context.Background(), "debug", debug) //nolint:staticcheck
}

func configureLogging(ctx context.Context) {
	if ctx.Value("debug").(bool) {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("debug mode turned on")
	}
}
